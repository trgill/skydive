/*
 * Copyright (C) 2019 Red Hat, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy ofthe License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specificlanguage governing permissions and
 * limitations under the License.
 *
 */

package blockdev

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/docker/docker/api/types/events"
	"github.com/skydive-project/skydive/common"
	"github.com/skydive-project/skydive/graffiti/graph"
	"github.com/skydive-project/skydive/probe"
	"github.com/skydive-project/skydive/topology"
	tp "github.com/skydive-project/skydive/topology/probes"
	ns "github.com/skydive-project/skydive/topology/probes/netns"
	libudev "github.com/trgill/go-libudev"
	udev "github.com/trgill/go-libudev"
)

type blockdevInfo struct {
	ID   string
	Node *graph.Node
}

// ProbeHandler describes a Docker topology graph that enhance the graph
type ProbeHandler struct {
	common.RWMutex
	*ns.ProbeHandler
	state       int64
	blockdevMap map[string]blockdevInfo
	wg          sync.WaitGroup
	connected   atomic.Value
}

func (p *ProbeHandler) blockdevNamespace(ID string) string {
	return fmt.Sprintf("/sys/block/%s", ID)
}

func (p *ProbeHandler) registerBlockdev(blockdev *udev.Device) {
	p.Lock()
	defer p.Unlock()

	if _, ok := p.blockdevMap[blockdev.Devnode()]; ok {
		return
	}

	p.Ctx.Graph.Lock()
	defer p.Ctx.Graph.Unlock()

	blockdevMetadata := Metadata{
		BlockDevID:    blockdev.Devnode(),
		BlockDevName:  blockdev.Devpath(),
		Sysname:       blockdev.Sysname(),
		Syspath:       blockdev.Syspath(),
		Devpath:       blockdev.Devpath(),
		Devnode:       blockdev.Devnode(),
		Subsystem:     blockdev.Subsystem(),
		Devtype:       blockdev.Devtype(),
		Sysnum:        blockdev.Sysnum(),
		IsInitialized: blockdev.IsInitialized(),
		Driver:        blockdev.Driver(),
	}

	node := p.Ctx.Graph.LookupFirstNode(graph.Metadata{"BlockDevID": blockdev.Devnode()})

	if node == nil {
		metadata := graph.Metadata{
			"Type":       "blockdev",
			"Name":       blockdev.Devnode(),
			"Manager":    "blockdev",
			"BlockDevID": blockdev.Devnode(),
			"Blockdev":   blockdevMetadata,
		}
		var err error
		if node, err = p.Ctx.Graph.NewNode(graph.GenID(), metadata); err != nil {
			p.Ctx.Logger.Error(err)
			return
		}

		parent := blockdev.Parent()
		if blockdev.Parent() != nil {
			parentNode := p.Ctx.Graph.LookupFirstNode(graph.Metadata{"BlockDevID": parent.Devnode()})
			if parentNode != nil {
				topology.AddOwnershipLink(p.Ctx.Graph, parentNode, node, nil)
			} else {
				topology.AddOwnershipLink(p.Ctx.Graph, p.Ctx.RootNode, node, nil)
			}
		}
	}

	p.blockdevMap[blockdev.Devnode()] = blockdevInfo{
		ID:   blockdev.Devnode(),
		Node: node,
	}
}

func (p *ProbeHandler) unregisterContainer(id string) {
	p.Lock()
	defer p.Unlock()

	infos, ok := p.blockdevMap[id]
	if !ok {
		return
	}

	p.Ctx.Graph.Lock()
	if err := p.Ctx.Graph.DelNode(infos.Node); err != nil {
		p.Ctx.Graph.Unlock()
		p.Ctx.Logger.Error(err)
		return
	}
	p.Ctx.Graph.Unlock()

	// namespace := p.containerNamespace(infos.ID)
	// p.Ctx.Logger.Debugf("Stop listening for namespace %s with PID %d", namespace, infos.ID)
	// p.Unregister(namespace)

	delete(p.blockdevMap, id)
}

func (p *ProbeHandler) handleBlockdevEvent(event *events.Message) {

}

func (p *ProbeHandler) connect() error {
	p.Ctx.Logger.Debugf("Connecting to udev")
	// Create Udev and Enumerate
	u := libudev.Udev{}
	e := u.NewEnumerate()

	// Add some FilterAddMatchSubsystemDevtype
	e.AddMatchSubsystem("block")
	//e.AddMatchIsInitialized()
	devices, _ := e.Devices()
	for i := range devices {
		device := devices[i]
		fmt.Println(device.Devnode())
		fmt.Println(device.Devlinks())
		fmt.Println(device.Properties())
		p.registerBlockdev(device)
	}

	return nil
}

// Start the probe
func (p *ProbeHandler) Start() {

	if !atomic.CompareAndSwapInt64(&p.state, common.StoppedState, common.RunningState) {
		return
	}

	go func() {
		for {
			state := atomic.LoadInt64(&p.state)
			if state == common.StoppingState || state == common.StoppedState {
				break
			}

			if p.connect() != nil {
				time.Sleep(time.Duration(24000) * time.Second)
			}

			p.wg.Wait()
		}
	}()
}

// Stop the probe
func (p *ProbeHandler) Stop() {
	if !atomic.CompareAndSwapInt64(&p.state, common.RunningState, common.StoppingState) {
		return
	}

	if p.connected.Load() == true {
		// p.cancel()
		p.wg.Wait()
	}

	atomic.StoreInt64(&p.state, common.StoppedState)
}

// Init initializes a new topology Docker probe
func (p *ProbeHandler) Init(ctx tp.Context, bundle *probe.Bundle) (probe.Handler, error) {
	nsHandler := bundle.GetHandler("netns")
	if nsHandler == nil {
		return nil, errors.New("unable to find the netns handler")
	}

	netnsRunPath := ctx.Config.GetString("agent.topology.docker.netns.run_path")

	p.ProbeHandler = nsHandler.(*ns.ProbeHandler)
	p.blockdevMap = make(map[string]blockdevInfo)
	p.state = common.StoppedState

	if netnsRunPath != "" {
		p.Exclude(netnsRunPath + "/default")
		p.Watch(netnsRunPath)
	}

	return p, nil
}

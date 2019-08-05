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
	"fmt"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/docker/docker/api/types/events"
	dm "github.com/docker/docker/pkg/devicemapper"
	"github.com/skydive-project/skydive/common"
	"github.com/skydive-project/skydive/graffiti/graph"
	"github.com/skydive-project/skydive/probe"
	"github.com/skydive-project/skydive/topology"
	tp "github.com/skydive-project/skydive/topology/probes"
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
	Ctx         tp.Context
	state       int64
	url         string
	blockdevMap map[string]blockdevInfo
	wg          sync.WaitGroup
	connected   atomic.Value
	Groups      map[string]*graph.Node
}

// Group describes blockdev groupings
type Group struct {
	GroupID uint
	Type    string
	Name    string
	UUID    string
}

func (p *ProbeHandler) blockdevNamespace(ID string) string {
	return fmt.Sprintf("/sys/block/%s", ID)
}

func (p *ProbeHandler) getLinks(blockdev *udev.Device, lines []dm.TargetLine) (links []string, err error) {

	//fmt.Println("lines:", lines)

	for _, line := range lines {
		paramsSplit := strings.Split(line.Params, " ")
		switch line.TargetType {
		case "linear":
			links = append(links, paramsSplit[0])
			//fmt.Println("links = ", links)
		case "raid":
			var index int64
			var raidDevs int64
			index = 1
			//fmt.Println("raidType:", paramsSplit[0], " #raid params: ", paramsSplit[1])
			raidParams, err := strconv.ParseInt(paramsSplit[index], 10, 64)
			if err != nil {
				return nil, err
			}
			index = index + raidParams + 1
			raidDevs, err = strconv.ParseInt(paramsSplit[index], 10, 64)
			index++
			var i int64

			for i = index; i < index+raidDevs*2; i++ {
				links = append(links, paramsSplit[i])
			}

		case "error":
			//fmt.Println("error:", line.Params)
		}
	}
	return links, nil
}
func (p *ProbeHandler) isMDTarget(blockdev *udev.Device) (isDM bool) {
	if blockdev.Properties()["UDISKS_MD_DEVICES"] != "" {
		return true
	}
	return false
}

func (p *ProbeHandler) isPartition(blockdev *udev.Device) (isDM bool) {
	if blockdev.Devtype() == "partition" {
		return true
	}
	return false
}

func (p *ProbeHandler) isDMTarget(blockdev *udev.Device) (isDM bool) {
	str := blockdev.Properties()["DM_ACTIVATION"]

	if str == "" {
		str = blockdev.Properties()["DM_NAME"]
	}
	if str == "" {
		str = blockdev.Properties()["DM_LV_NAME"]
	}

	return str != ""
}

// getGroup returns a Group struct
func (p *ProbeHandler) getGroup(groupID uint, groupType string, name string, uuid string) Group {

	return Group{
		GroupID: groupID,
		Type:    groupType,
		Name:    name,
		UUID:    uuid,
	}
}

// addGroup adds a group to the graph and links it to the bridge
func (p *ProbeHandler) addGroup(group Group) {
	fmt.Println("New group added", group.UUID)
	g := p.Ctx.Graph
	g.Lock()
	defer g.Unlock()
	metadata := graph.Metadata{
		"Type":      "blockdev",
		"GroupId":   group.GroupID,
		"GroupType": group.Type,
		"UUID":      group.UUID,
	}
	groupNode, err := g.NewNode(graph.GenID(), metadata)
	if err != nil {
		p.Ctx.Logger.Error(err)
		return
	}
	p.Groups["dm"] = groupNode
	// if _, err := topology.AddOwnershipLink(g, p.Ctx.RootNode, groupNode, nil); err != nil {
	// 	p.Ctx.Logger.Error(err)
	// }
}

func (p *ProbeHandler) getName(blockdev *udev.Device) (name string) {

	if p.isDMTarget(blockdev) {
		name = blockdev.Properties()["DM_NAME"]
	} else {
		name = blockdev.Devnode()
	}
	//fmt.Println("devnode", blockdev.Devnode(), " name = ", name)
	return name
}

func (p *ProbeHandler) getMetaData(blockdev *udev.Device) (metadata graph.Metadata, linkedDevs []string) {
	var blockdevMetadata Metadata
	majorMinor := blockdev.Properties()["MAJOR"] + ":" + blockdev.Properties()["MINOR"]
	name := p.getName(blockdev)

	blockdevMetadata = Metadata{
		BlockDevID:    name,
		MajorMinor:    majorMinor,
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

	if blockdev.Devtype() == "partition" {
		fmt.Println("Partition Dependency = ", blockdev.Properties()["ID_PART_ENTRY_DISK"])
		linkedDevs = append(linkedDevs, blockdev.Properties()["ID_PART_ENTRY_DISK"])

	} else if p.isDMTarget(blockdev) {
		lines, err := dm.GetTable(blockdev.Devnode())
		if err == nil {
			fmt.Println("Table = ", lines)
			linkedDevs, _ = p.getLinks(blockdev, lines)
		}
	} else if p.isMDTarget(blockdev) {
		for k, v := range blockdev.Properties() {
			if strings.HasPrefix(k, "MD_DEVICE") && strings.HasSuffix(k, "_DEV") {
				fmt.Println("MD dev = ", v)
				linkedDevs = append(linkedDevs, v)
			}
		}

	}

	metadata = graph.Metadata{
		"Type":       "blockdev",
		"Name":       name,
		"MajorMinor": majorMinor,
		"Manager":    "blockdev",
		"BlockDevID": blockdev.Devnode(),
		"Blockdev":   blockdevMetadata,
		"LinkedDevs": linkedDevs,
	}
	fmt.Println("linkedDevs", linkedDevs)
	return metadata, linkedDevs
}

func (p *ProbeHandler) registerBlockdev(blockdev *udev.Device, dmGroup Group) {
	p.Lock()
	defer p.Unlock()
	var linkedDevs []string
	var graphMetaData graph.Metadata

	if _, ok := p.blockdevMap[blockdev.Devnode()]; ok {
		return
	}

	p.Ctx.Graph.Lock()
	defer p.Ctx.Graph.Unlock()
	name := p.getName(blockdev)
	fmt.Println("Lookup Name : ", name)
	node := p.Ctx.Graph.LookupFirstNode(graph.Metadata{"BlockDevID": name})

	if node == nil {

		graphMetaData, linkedDevs = p.getMetaData(blockdev)

		var err error
		if node, err = p.Ctx.Graph.NewNode(graph.GenID(), graphMetaData); err != nil {
			p.Ctx.Logger.Error(err)
			return
		}
	}

	if blockdev.Devtype() == "partition" {
		parent := blockdev.Parent()
		if parent != nil {
			fmt.Println("parent = ", parent.Devnode())
			parentNode := p.Ctx.Graph.LookupFirstNode(graph.Metadata{"BlockDevID": parent.Devnode()})
			if parentNode != nil {
				topology.AddLink(p.Ctx.Graph, parentNode, node, "connected", nil)
			} else {
				topology.AddLink(p.Ctx.Graph, node, p.Ctx.RootNode, "connected", nil)
			}
		}
	} else if p.isDMTarget(blockdev) {
		for _, connected := range linkedDevs {
			fmt.Println("connected dev = ", connected)
			connectedNode := p.Ctx.Graph.LookupFirstNode(graph.Metadata{"MajorMinor": connected})
			if connectedNode != nil {
				topology.AddLink(p.Ctx.Graph, node, connectedNode, "connected", nil)
				topology.AddOwnershipLink(p.Ctx.Graph, p.Groups["dm"], node, nil)
			}
		}
	} else if blockdev.Properties()["UDISKS_MD_DEVICES"] != "" {
		for _, connected := range linkedDevs {
			fmt.Println("connected dev = ", connected)
			connectedNode := p.Ctx.Graph.LookupFirstNode(graph.Metadata{"BlockDevID": connected})
			if connectedNode != nil {
				topology.AddLink(p.Ctx.Graph, node, connectedNode, "connected", nil)
			}
		}
	} else {
		topology.AddOwnershipLink(p.Ctx.Graph, p.Ctx.RootNode, node, nil)
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

func (p *ProbeHandler) connect(dmGroup Group) error {
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
		fmt.Println("Start: ", device.Devnode())
		//fmt.Println(device.Devlinks())
		fmt.Println(device.Properties())
		//fmt.Println(device.Sysattrs())

		p.registerBlockdev(device, dmGroup)

	}

	return nil
}

// Start the probe
func (p *ProbeHandler) Start() {

	if !atomic.CompareAndSwapInt64(&p.state, common.StoppedState, common.RunningState) {
		return
	}
	p.Groups = make(map[string]*graph.Node)
	dmGroup := p.getGroup(1, "device mapper", "device mapper targets", "2c35690e-b422-4c4b-80bd-e34ea40b4e8e")
	p.addGroup(dmGroup)
	go func() {
		for {
			state := atomic.LoadInt64(&p.state)
			if state == common.StoppingState || state == common.StoppedState {
				break
			}
			// TODO: connect() should return a result
			if p.connect(dmGroup) == nil {
				time.Sleep(5 * time.Minute)
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

// Init initializes a new topology blockdev probe
func (p *ProbeHandler) Init(ctx tp.Context, bundle *probe.Bundle) (probe.Handler, error) {

	blockdevURL := ctx.Config.GetString("agent.topology.blockdev.url")

	p.Ctx = ctx
	p.url = blockdevURL
	p.blockdevMap = make(map[string]blockdevInfo)
	p.state = common.StoppedState

	return p, nil
}

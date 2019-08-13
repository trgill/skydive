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
	"io"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/coreos/go-systemd/sdjournal"
	"github.com/docker/docker/api/types/events"
	"github.com/skydive-project/skydive/common"
	"github.com/skydive-project/skydive/graffiti/graph"
	"github.com/skydive-project/skydive/probe"
	"github.com/skydive-project/skydive/topology"
	tp "github.com/skydive-project/skydive/topology/probes"
	ns "github.com/skydive-project/skydive/topology/probes/netns"
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

func (p *ProbeHandler) followJournal() {

	p.Ctx.Logger.Debugf("Connecting to journal")

	r, err := sdjournal.NewJournalReader(sdjournal.JournalReaderConfig{
		Since: time.Duration(-15) * time.Second,
		Matches: []sdjournal.Match{
			{
				Field: sdjournal.SD_JOURNAL_FIELD_SYSTEMD_UNIT,
				Value: "NetworkManager.service",
			},
		},
	})

	if err != nil {
		p.Ctx.Logger.Fatalf("Error opening journal: %s", err)
	}

	if r == nil {
		p.Ctx.Logger.Fatal("Got a nil reader")
	}
	defer r.Close()

	timeout := time.Duration(50) * time.Second
	if err = r.Follow(time.After(timeout), os.Stdout); err != sdjournal.ErrExpired {
		p.Ctx.Logger.Fatalf("Error during follow: %s", err)
	}

}
func (p *ProbeHandler) getValue(j *sdjournal.Journal, searchValue string) string {

	value, err := j.GetDataValue(searchValue)

	if err != nil {
		return ""
	}
	return value
}

func (p *ProbeHandler) getMetaData(j *sdjournal.Journal) (*Metadata, error) {

	var key string

	devName, err := j.GetDataValue("DEVNAME")
	if err == nil && len(devName) > 0 {
		key = devName
	}
	dmUUID, err := j.GetDataValue("DM_UUID")
	if err == nil && len(dmUUID) > 0 {
		key = dmUUID
	}
	storageID, err := j.GetDataValue("PERSISTENT_STORAGE_ID")
	if err == nil && len(storageID) > 0 {
		key = storageID
	}
	storageLabel, err := j.GetDataValue("PERSISTENT_STORAGE_LABEL")
	if err == nil && len(storageLabel) > 0 {
		key = storageLabel
	}
	partLabel, err := j.GetDataValue("PERSISTENT_STORAGE_PARTLABEL")
	if err == nil && len(partLabel) > 0 {
		key = partLabel
	}
	partUUID, err := j.GetDataValue("PERSISTENT_STORAGE_PARTUUID")
	if err == nil && len(partUUID) > 0 {
		key = partUUID
	}
	storagePath, err := j.GetDataValue("PERSISTENT_STORAGE_PATH")
	if err == nil && len(storagePath) > 0 {
		key = storagePath
	}
	storageUUID, err := j.GetDataValue("PERSISTENT_STORAGE_UUID")
	if err == nil && len(storageUUID) > 0 {
		key = storageUUID
	}

	metadata := Metadata{
		key:                                   key,
		storageID:                             storageID,
		storageLabel:                          storageLabel,
		partLabel:                             partLabel,
		partUUID:                              partUUID,
		storagePath:                           storagePath,
		storageUUID:                           storageUUID,
		DEVTYPE:                               p.getValue(j, "DEVTYPE"),
		PWD:                                   p.getValue(j, "PWD"),
		SHLVL:                                 p.getValue(j, "SHLVL"),
		SUBSYSTEM:                             p.getValue(j, "SUBSYSTEM"),
		TAGS:                                  p.getValue(j, "TAGS"),
		UDEVLOG_QUEUE_DISCARD_GRANULARIT:      p.getValue(j, "UDEVLOG_QUEUE_DISCARD_GRANULARIT"),
		UDEVLOG_QUEUE_DISCARD_MAX_BYTES:       p.getValue(j, "UDEVLOG_QUEUE_DISCARD_MAX_BYTES"),
		UDEVLOG_QUEUE_DISCARD_ZEROES_DATA:     p.getValue(j, "UDEVLOG_QUEUE_DISCARD_ZEROES_DATA"),
		UDEVLOG_QUEUE_LOGICAL_BLOCK_SIZE:      p.getValue(j, "UDEVLOG_QUEUE_LOGICAL_BLOCK_SIZE"),
		UDEVLOG_QUEUE_MINIMUM_IO_SIZE:         p.getValue(j, "UDEVLOG_QUEUE_MINIMUM_IO_SIZE"),
		UDEVLOG_QUEUE_OPTIMAL_IO_SIZE:         p.getValue(j, "UDEVLOG_QUEUE_OPTIMAL_IO_SIZE"),
		UDEVLOG_QUEUE_PHYSICAL_BLOCK_SIZE:     p.getValue(j, "UDEVLOG_QUEUE_PHYSICAL_BLOCK_SIZE"),
		UDEVLOG_QUEUE_READ_AHEAD_KB:           p.getValue(j, "UDEVLOG_QUEUE_READ_AHEAD_KB"),
		UDEVLOG_QUEUE_ROTATIONAL:              p.getValue(j, "UDEVLOG_QUEUE_ROTATIONAL"),
		UDEVLOG_BLOCK_ALIGNMENT_OFFSET:        p.getValue(j, "UDEVLOG_BLOCK_ALIGNMENT_OFFSET"),
		UDEVLOG_BLOCK_DISCARD_ALIGNMENT:       p.getValue(j, "UDEVLOG_BLOCK_DISCARD_ALIGNMENT"),
		UDEVLOG_BLOCK_REMOVABLE:               p.getValue(j, "UDEVLOG_BLOCK_REMOVABLE"),
		UDEVLOG_BLOCK_RO:                      p.getValue(j, "UDEVLOG_BLOCK_RO"),
		UDEVLOG_LSBLK_MOUNTPOINT:              p.getValue(j, "UDEVLOG_LSBLK_MOUNTPOINT"),
		UDEVLOG_LSBLK_LABEL:                   p.getValue(j, "UDEVLOG_LSBLK_LABEL"),
		UDEVLOG_LSBLK_STATE:                   p.getValue(j, "UDEVLOG_LSBLK_STATE"),
		UDEVLOG_LSBLK_OWNER:                   p.getValue(j, "UDEVLOG_LSBLK_OWNER"),
		UDEVLOG_LSBLK_GROUP:                   p.getValue(j, "UDEVLOG_LSBLK_GROUP"),
		_COMM:                                 p.getValue(j, "_COMM"),
		_SELINUX_CONTEXT:                      p.getValue(j, "_SELINUX_CONTEXT"),
		_SYSTEMD_CGROUP:                       p.getValue(j, "_SYSTEMD_CGROUP"),
		_SYSTEMD_UNIT:                         p.getValue(j, "_SYSTEMD_UNIT"),
		_SYSTEMD_INVOCATION_ID:                p.getValue(j, "_SYSTEMD_INVOCATION_ID"),
		DM_ACTIVATION:                         p.getValue(j, "DM_ACTIVATION"),
		DM_LV_LAYER:                           p.getValue(j, "DM_LV_LAYER"),
		DM_NOSCAN:                             p.getValue(j, "DM_NOSCAN"),
		DM_SUBSYSTEM_UDEV_FLAG0:               p.getValue(j, "DM_SUBSYSTEM_UDEV_FLAG0"),
		DM_SUSPENDED:                          p.getValue(j, "DM_SUSPENDED"),
		DM_UDEV_DISABLE_LIBRARY_FALLBACK_FLAG: p.getValue(j, "DM_UDEV_DISABLE_LIBRARY_FALLBACK_FLAG"),
		DM_UDEV_DISABLE_OTHER_RULES_FLAG:      p.getValue(j, "DM_UDEV_DISABLE_OTHER_RULES_FLAG"),
		DM_UDEV_PRIMARY_SOURCE_FLAG:           p.getValue(j, "DM_UDEV_PRIMARY_SOURCE_FLAG"),
		DM_UDEV_RULES_VSN:                     p.getValue(j, "DM_UDEV_RULES_VSN"),
		DM_VG_NAME:                            p.getValue(j, "DM_VG_NAME"),
		MAJOR:                                 p.getValue(j, "MAJOR"),
		UDEVLOG_QUEUE_NR_REQUESTS:             p.getValue(j, "UDEVLOG_QUEUE_NR_REQUESTS"),
		UDEVLOG_QUEUE_SCHEDULER:               p.getValue(j, "UDEVLOG_QUEUE_SCHEDULER"),
		UDEVLOG_BLOCK_CAPABILITY:              p.getValue(j, "UDEVLOG_BLOCK_CAPABILITY"),
		UDEVLOG_BLOCK_EXT_RANGE:               p.getValue(j, "UDEVLOG_BLOCK_EXT_RANGE"),
		UDEVLOG_LSBLK_MODE:                    p.getValue(j, "UDEVLOG_LSBLK_MODE"),
		UDEVLOG_LSBLK_TYPE:                    p.getValue(j, "UDEVLOG_LSBLK_TYPE"),
		UDEVLOG_DM_READ_AHEAD:                 p.getValue(j, "UDEVLOG_DM_READ_AHEAD"),
		UDEVLOG_DM_ATTR:                       p.getValue(j, "UDEVLOG_DM_ATTR"),
		UDEVLOG_DM_TABLES_LOADED:              p.getValue(j, "UDEVLOG_DM_TABLES_LOADED"),
		UDEVLOG_DM_READONLY:                   p.getValue(j, "UDEVLOG_DM_READONLY"),
		UDEVLOG_DM_SEGMENTS:                   p.getValue(j, "UDEVLOG_DM_SEGMENTS"),
		UDEVLOG_DM_EVENTS:                     p.getValue(j, "UDEVLOG_DM_EVENTS"),
		UDEVLOG_DM_SUBSYSTEM:                  p.getValue(j, "UDEVLOG_DM_SUBSYSTEM"),
		UDEVLOG_DM_LV_LAYER:                   p.getValue(j, "UDEVLOG_DM_LV_LAYER"),
		UDEVLOG_DM_TABLE_INACTIVE:             p.getValue(j, "UDEVLOG_DM_TABLE_INACTIVE"),
		UDEVLOG_BLOCK_SIZE:                    p.getValue(j, "UDEVLOG_BLOCK_SIZE"),
		DEVLINKS:                              p.getValue(j, "DEVLINKS"),
		DEVNAME:                               p.getValue(j, "DEVNAME"),
		DEVPATH:                               p.getValue(j, "DEVPATH"),
		DM_COOKIE:                             p.getValue(j, "DM_COOKIE"),
		DM_DISABLE_OTHER_RULES_FLAG_OLD:       p.getValue(j, "DM_DISABLE_OTHER_RULES_FLAG_OLD"),
		DM_LV_NAME:                            p.getValue(j, "DM_LV_NAME"),
		DM_NAME:                               p.getValue(j, "DM_NAME"),
		DM_UUID:                               p.getValue(j, "DM_UUID"),
		MINOR:                                 p.getValue(j, "MINOR"),
		SEQNUM:                                p.getValue(j, "SEQNUM"),
		USEC_INITIALIZED:                      p.getValue(j, "USEC_INITIALIZED"),
		PERSISTENT_STORAGE_ID:                 p.getValue(j, "PERSISTENT_STORAGE_ID"),
		UDEVLOG_DM_OPEN:                       p.getValue(j, "UDEVLOG_DM_OPEN"),
		UDEVLOG_DM_DEVICE_COUNT:               p.getValue(j, "UDEVLOG_DM_DEVICE_COUNT"),
		UDEVLOG_DM_DEVS_USED:                  p.getValue(j, "UDEVLOG_DM_DEVS_USED"),
		UDEVLOG_DM_DEVNOS_USED:                p.getValue(j, "UDEVLOG_DM_DEVNOS_USED"),
		UDEVLOG_DM_BLKDEVS_USED:               p.getValue(j, "UDEVLOG_DM_BLKDEVS_USED"),
		UDEVLOG_DM_DEVICE_REF_COUNT:           p.getValue(j, "UDEVLOG_DM_DEVICE_REF_COUNT"),
		UDEVLOG_DM_NAMES_USING_DEV:            p.getValue(j, "UDEVLOG_DM_NAMES_USING_DEV"),
		UDEVLOG_DM_DEVNOS_USING_DEV:           p.getValue(j, "UDEVLOG_DM_DEVNOS_USING_DEV"),
		UDEVLOG_DM_TABLE_LIVE:                 p.getValue(j, "UDEVLOG_DM_TABLE_LIVE"),
	}
	return &metadata, nil
}

func (p *ProbeHandler) addGraphLinkForID(persistentID string, linkType string, node *graph.Node) {

	p.Ctx.Logger.Info("Connecting to ID:", persistentID)

	id := strings.Split(persistentID, "-part")[0]
	child := p.Ctx.Graph.LookupFirstNode(graph.Metadata{"BlockDevID": id})

	if child != nil {
		topology.AddLink(p.Ctx.Graph, child, node, linkType, nil)
	}

}

func (p *ProbeHandler) addGraphLinkForList(parentBlockDevs string, childBlockDevs string, node *graph.Node) {
	parentBlockDevsArray := strings.Split(strings.Trim(parentBlockDevs, "'"), ",")
	for _, majorMinor := range parentBlockDevsArray {
		majorMinorSplit := strings.Split(majorMinor, ":")

		p.Ctx.Logger.Info("Connecting to MAJOR:", majorMinorSplit[0], "MINOR:", majorMinorSplit[1])
		child := p.Ctx.Graph.LookupFirstNode(graph.Metadata{"MAJOR": majorMinorSplit[0], "MINOR": majorMinorSplit[1]})

		if child != nil {
			topology.AddLink(p.Ctx.Graph, child, node, "uses", nil)
		}
	}

	childBlockDevsArray := strings.Split(strings.Trim(childBlockDevs, "'"), ",")
	for _, dmName := range childBlockDevsArray {

		p.Ctx.Logger.Info("Connecting to DM_NAME:", dmName)
		child := p.Ctx.Graph.LookupFirstNode(graph.Metadata{"DM_NAME": dmName})

		if child != nil {
			topology.AddLink(p.Ctx.Graph, child, node, "used-by", nil)
		}
	}
}

func (p *ProbeHandler) getName(blockdevMetadata *Metadata) string {

	if len(blockdevMetadata.DM_NAME) > 0 {
		return blockdevMetadata.DM_NAME
	}

	return blockdevMetadata.DEVNAME
}

func (p *ProbeHandler) registerBlockdev(j *sdjournal.Journal) {
	p.Lock()
	defer p.Unlock()

	p.Ctx.Graph.Lock()
	defer p.Ctx.Graph.Unlock()

	blockdevMetadata, err := p.getMetaData(j)

	if err != nil {
		p.Ctx.Logger.Error(err)
		return
	}

	key := blockdevMetadata.key
	node := p.Ctx.Graph.LookupFirstNode(graph.Metadata{"BlockDevID": key})

	p.Ctx.Logger.Info("type: ", blockdevMetadata.DEVTYPE, "name", blockdevMetadata.DEVNAME, "has key : ", key)

	if node == nil {
		metadata := graph.Metadata{
			"Type":       "blockdev",
			"Name":       p.getName(blockdevMetadata),
			"Manager":    "blockdev",
			"BlockDevID": key,
			"Blockdev":   blockdevMetadata,
			"MAJOR":      blockdevMetadata.MAJOR,
			"MINOR":      blockdevMetadata.MINOR,
			"DM_NAME":    blockdevMetadata.DM_NAME,
		}
		var err error
		if node, err = p.Ctx.Graph.NewNode(graph.GenID(), metadata); err != nil {
			p.Ctx.Logger.Error(err)
			return
		}

		if len(blockdevMetadata.UDEVLOG_DM_DEVNOS_USED) > 0 || len(blockdevMetadata.UDEVLOG_DM_DEVNOS_USING_DEV) > 0 {
			p.addGraphLinkForList(blockdevMetadata.UDEVLOG_DM_DEVNOS_USED, blockdevMetadata.UDEVLOG_DM_NAMES_USING_DEV, node)
		} else if blockdevMetadata.DEVTYPE == "partition" && len(blockdevMetadata.PERSISTENT_STORAGE_ID) > 0 {
			p.addGraphLinkForID(blockdevMetadata.PERSISTENT_STORAGE_ID, "uses", node)
		} else {
			topology.AddOwnershipLink(p.Ctx.Graph, p.Ctx.RootNode, node, nil)
		}
	}

	p.blockdevMap[key] = blockdevInfo{
		ID:   key,
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

// matchField := "PERSISTENT_STORAGE_ID"
// matchValue := "dm-name-mirror_vg-mirror_lv"

func (p *ProbeHandler) connect() error {

	j, err := sdjournal.NewJournal()
	if err != nil {
		p.Ctx.Logger.Info("Error opening journal: %s", err)
		return err
	}
	p.Ctx.Logger.Info(j)

	if j == nil {
		p.Ctx.Logger.Info("Got a nil journal")
		return err
	}
	defer j.Close()

	if err = j.SeekHead(); err != nil {
		p.Ctx.Logger.Info(err)
		return err
	}

	matchField := "DEVTYPE"
	m := sdjournal.Match{Field: matchField, Value: "disk"}
	if err = j.AddMatch(m.String()); err != nil {
		p.Ctx.Logger.Info("Error adding matches to journal: %s", err)
		return err
	}
	matchField = "DEVTYPE"
	m = sdjournal.Match{Field: matchField, Value: "partition"}
	if err = j.AddMatch(m.String()); err != nil {
		p.Ctx.Logger.Info("Error adding matches to journal: %s", err)
		return err
	}
	matchField = "SUBSYSTEM"
	m = sdjournal.Match{Field: matchField, Value: "block"}
	if err = j.AddMatch(m.String()); err != nil {
		p.Ctx.Logger.Info("Error adding matches to journal: %s", err)
	}

	for {

		n, err := j.Next()

		if err != nil {
			return fmt.Errorf("Error reading to journal: %s", err)
		}

		if n == 0 {
			return fmt.Errorf("Error reading to journal: %s", io.EOF)
		}

		p.registerBlockdev(j)
	}

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

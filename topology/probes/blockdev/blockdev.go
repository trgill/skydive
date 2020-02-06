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
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"path"
	"sync"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/skydive-project/skydive/common"
	"github.com/skydive-project/skydive/graffiti/graph"
	"github.com/skydive-project/skydive/probe"
	"github.com/skydive-project/skydive/topology"
	"github.com/skydive-project/skydive/topology/probes"
	tp "github.com/skydive-project/skydive/topology/probes"
)

const blockdevGroupName string = "block devices"

// Types returned from lsblk JSON
const multipathType string = "mpath"
const lvmType string = "lvm"
const blockGroupType string = "cluster"

// Types for skydive nodes
const lvmNodeType string = "blockdevlvm"
const blockdevNodeType string = "blockdev"
const leafNodeType string = "blockdevleaf"

const managerType string = "blockdev"

const defaultIostatInterval int = 300

type statistics struct {
	Metrics []IOMetric `mapstructure:"disk"`
}

type hostdata struct {
	Nodename   string       `mapstructure:"nodename"`
	Sysname    string       `mapstructure:"sysname"`
	Release    string       `mapstructure:"release"`
	Machine    string       `mapstructure:"machine"`
	Cpucount   uint32       `mapstructure:"number-of-cpus"`
	Date       string       `mapstructure:"date"`
	Statistics []statistics `mapstructure:"statistics"`
}

type hosts struct {
	Hosts []hostdata `mapstructure:"hosts"`
}
type sysstat struct {
	Sysstat hosts `mapstructure:"sysstat"`
}

// BlockDeviceAttrs use to hold IO stats for a block device
type BlockDeviceAttrs struct {
	tps        int64
	kBReadPerS int64
	kBWrtnPerS int64
	kBRead     int64
	kBWrtn     int64
}

type VolumeGroup struct {
	Name       string `json:"vg_name,omitempty"`
	Attr       string `json:"vg_attr,omitempty"`
	ExtentSize string `json:"vg_extent_size,omitempty"`
	PVCount    string `json:"pv_count,omitempty"`
	LVCount    string `json:"lv_count,omitempty"`
	SnanpCount string `json:"snap_count,omitempty"`
	Size       string `json:"vg_size,omitempty"`
	Free       string `json:"vg_free,omitempty"`
	UUID       string `json:"vg_uuid,omitempty"`
	Profile    string `json:"vg_profile,omitempty"`
}

type VGReports struct {
	LVList []VolumeGroup `json:"vg,omitempty"`
}

type VGReport struct {
	Report []VGReports `json:"report,omitempty"`
}

type PhysicalVolume struct {
	Name    string `json:"pv_name,omitempty"`
	VG      string `json:"vg_name",omitempty"`
	Fmt     string `json:"pv_fmt",omitempty"`
	Attr    string `json:"pv_attr",omitempty"`
	Size    string `json:"pv_size",omitempty"`
	Free    string `json:"pv_free",omitempty"`
	Major   string `json:"pv_major",omitempty"`
	Minor   string `json:"pv_minor",omitempty"`
	DevSize string `json:"dev_size",omitempty"`
	UUID    string `json:"pv_uuid",omitempty"`
}

type PVReports struct {
	LVList []PhysicalVolume `json:"pv,omitempty"`
}

type PVReport struct {
	Report []PVReports `json:"report,omitempty"`
}

type LogicalVolume struct {
	Name            string `json:"lv_name,omitempty"`
	Path            string `json:"lv_path,omitempty"`
	DMPath          string `json:"lv_dm_path,omitempty"`
	Layout          string `json:"lv_layout,omitempty"`
	VG              string `json:"vg_name,omitempty"`
	Attr            string `json:"lv_attr,omitempty"`
	Size            string `json:"lv_size,omitempty"`
	Pool            string `json:"pool_lv,omitempty"`
	Origin          string `json:"origin,omitempty"`
	DataPercent     string `json:"data_percent,omitempty"`
	MetadataPercent string `json:"metadata_percent,omitempty"`
	MovePV          string `json:"move_pv,omitempty"`
	MirrorLog       string `json:"mirror_log,omitempty"`
	CopyPercent     string `json:"copy_percent,omitempty"`
	ConvertLV       string `json:"convert_lv,omitempty"`
	KernelMajor     string `json:"lv_kernel_major,omitempty"`
	KernelMinor     string `json:"lv_kernel_minor,omitempty"`
	UUID            string `json:"lv_uuid,omitempty"`
}

type LVReports struct {
	LVList []LogicalVolume `json:"lv,omitempty"`
}

type LVReport struct {
	Report []LVReports `json:"report,omitempty"`
}

// BlockDevice used for JSON parsing
type BlockDevice struct {
	Children []BlockDevice `mapstructure:"children"`

	Alignment    int64  `mapstructure:"alignment"`
	DiscAln      int64  `mapstructure:"disc-aln"`
	DiscGran     string `mapstructure:"disc-gran"`
	DiscMax      string `mapstructure:"disc-max"`
	DiscZero     bool   `mapstructure:"disc-zero"`
	Fsavail      string `mapstructure:"fsavail"`
	Fssize       string `mapstructure:"fssize"`
	Fstype       string `mapstructure:"fstype"`
	FsusePercent string `mapstructure:"fsuse%"`
	Fsused       string `mapstructure:"fsused"`
	Group        string `mapstructure:"group"`
	Hctl         string `mapstructure:"hctl"`
	Hotplug      bool   `mapstructure:"hotplug"`
	Kname        string `mapstructure:"kname"`
	Label        string `mapstructure:"label"`
	LogSec       int64  `mapstructure:"log-sec"`
	MajMin       string `mapstructure:"maj:min"`
	MinIo        int64  `mapstructure:"min-io"`
	Mode         string `mapstructure:"mode"`
	Model        string `mapstructure:"model"`
	Mountpoint   string `mapstructure:"mountpoint"`
	Name         string `mapstructure:"name"`
	OptIo        int64  `mapstructure:"opt-io"`
	Owner        string `mapstructure:"owner"`
	Partflags    string `mapstructure:"partflags"`
	Partlabel    string `mapstructure:"partlabel"`
	Parttype     string `mapstructure:"parttype"`
	Partuuid     string `mapstructure:"partuuid"`
	Path         string `mapstructure:"path"`
	PhySec       int64  `mapstructure:"hpy-sec"`
	Pkname       string `mapstructure:"pkname"`
	Pttype       string `mapstructure:"pttype"`
	Ptuuid       string `mapstructure:"ptuuid"`
	Ra           int64  `mapstructure:"ra"`
	Rand         bool   `mapstructure:"rand"`
	Rev          string `mapstructure:"rev"`
	Rm           bool   `mapstructure:"rm"`
	Ro           bool   `mapstructure:"ro"`
	Rota         bool   `mapstructure:"rota"`
	RqSize       int64  `mapstructure:"rq-size"`
	Sched        string `mapstructure:"sched"`
	Serial       string `mapstructure:"serial"`
	Size         string `mapstructure:"size"`
	State        string `mapstructure:"state"`
	Subsystems   string `mapstructure:"subsystems"`
	Tran         string `mapstructure:"tran"`
	Type         string `mapstructure:"type"`
	UUID         string `mapstructure:"uuid"`
	Vendor       string `mapstructure:"vendor"`
	Wsame        string `mapstructure:"wsame"`
	WWN          string `mapstructure:"wwn"`
}

// Devices used for JSON parsing
type Devices struct {
	Blockdevices []BlockDevice `mapstructure:"blockdevices"`
}

type blockdevInfo struct {
	ID   string
	Name string
	Path string
	Node *graph.Node
}

// ProbeHandler describes a block device graph that enhances the graph
type ProbeHandler struct {
	common.RWMutex
	Ctx         tp.Context
	blockdevMap map[string]blockdevInfo
	wg          sync.WaitGroup
	Groups      map[string]*graph.Node
}

func (bd *BlockDevice) getPath() string {
	if bd.Path != "" {
		return bd.Path
	}
	if bd.Name != "" {
		return bd.Name
	}
	return bd.Kname
}

func (bd *BlockDevice) getID() string {
	if bd.WWN != "" {
		return bd.WWN
	}
	if bd.Serial != "" {
		return bd.Serial
	}
	return bd.Name
}

func (bd *BlockDevice) getName() string {
	if bd.Mountpoint != "" {
		return bd.Mountpoint
	}
	base := path.Base(bd.getPath())
	return base
}

func (p *ProbeHandler) addGroupByName(name string, ID string) *graph.Node {
	p.Lock()
	defer p.Unlock()
	if p.Groups[name] != nil {
		return p.Groups[name]
	}
	g := p.Ctx.Graph
	g.Lock()
	defer g.Unlock()
	metadata := graph.Metadata{
		"Name":    name,
		"ID":      ID,
		"Type":    blockGroupType,
		"Manager": managerType,
	}
	groupNode, err := g.NewNode(graph.GenID(), metadata)
	if err != nil {
		p.Ctx.Logger.Error(err)
		return nil
	}
	p.Groups[name] = groupNode

	return groupNode
}

// addGroup adds a group to the graph
func (p *ProbeHandler) addGroup(blockdev BlockDevice, ID string) *graph.Node {
	groupName := blockdev.getName()
	return p.addGroupByName(groupName, ID)
}

func (p *ProbeHandler) getPVMetaData(pv PhysicalVolume) (metadata graph.Metadata) {
	pvMetadata := PVMetadata{
		Name:    pv.Name,
		VG:      pv.VG,
		Fmt:     pv.Fmt,
		Attr:    pv.Attr,
		Size:    pv.Size,
		Free:    pv.Free,
		DevSize: pv.DevSize,
		Major:   pv.Major,
		Minor:   pv.Minor,
		UUID:    pv.UUID,
	}

	metadata = graph.Metadata{
		"PV":         pvMetadata,
		"Type":       lvmNodeType,
		"MajorMinor": pv.Major + ":" + pv.Minor,
	}

	return metadata
}
func (p *ProbeHandler) getVGMetaData(vg VolumeGroup) (metadata graph.Metadata) {
	vgMetadata := VGMetadata{
		Name:       vg.Name,
		Attr:       vg.Attr,
		ExtentSize: vg.ExtentSize,
		PVCount:    vg.PVCount,
		LVCount:    vg.LVCount,
		SnanpCount: vg.SnanpCount,
		Size:       vg.Size,
		Free:       vg.Free,
		UUID:       vg.UUID,
		Profile:    vg.Profile,
	}

	metadata = graph.Metadata{
		"VG": vgMetadata,
	}

	return metadata
}
func (p *ProbeHandler) getLVMetaData(lv LogicalVolume) (metadata graph.Metadata) {
	lvMetadata := LVMetadata{
		Name:            lv.Name,
		Path:            lv.Path,
		DMPath:          lv.DMPath,
		VG:              lv.VG,
		Attr:            lv.Attr,
		Size:            lv.Size,
		Pool:            lv.Pool,
		Origin:          lv.Origin,
		DataPercent:     lv.DataPercent,
		MetadataPercent: lv.MetadataPercent,
		MovePV:          lv.MovePV,
		MirrorLog:       lv.MirrorLog,
		CopyPercent:     lv.CopyPercent,
		ConvertLV:       lv.ConvertLV,
		KernelMajor:     lv.KernelMajor,
		KernelMinor:     lv.KernelMinor,
		UUID:            lv.UUID,
	}

	metadata = graph.Metadata{
		"LV":         lvMetadata,
		"MajorMinor": lv.KernelMajor + ":" + lv.KernelMinor,
		"Type":       lvmNodeType,
		"Path":       lv.Path,
		"DMPath":     lv.DMPath,
	}

	return metadata
}

func (p *ProbeHandler) getMetaData(blockdev BlockDevice, childCount int,
	parentWWN string) (metadata graph.Metadata) {
	var blockdevMetadata Metadata
	var nodeType string

	// The nodeType maps to the icon is used for display
	if childCount > 0 {
		if blockdev.Type == lvmType {
			nodeType = lvmNodeType
		} else {
			nodeType = blockdevNodeType
		}
	} else {
		nodeType = leafNodeType
	}

	// JSON for multipath devices doesn't include the WWN, so get it from
	// the parent for a multipathType
	var WWN string
	if blockdev.Type == multipathType {
		WWN = parentWWN
	} else {
		WWN = blockdev.WWN
	}

	blockdevMetadata = Metadata{
		Name:         blockdev.getName(),
		Alignment:    blockdev.Alignment,
		DiscAln:      blockdev.DiscAln,
		DiscGran:     blockdev.DiscGran,
		DiscMax:      blockdev.DiscMax,
		DiscZero:     blockdev.DiscZero,
		Fsavail:      blockdev.Fsavail,
		Fssize:       blockdev.Fssize,
		Fstype:       blockdev.Fstype,
		FsusePercent: blockdev.FsusePercent,
		Fsused:       blockdev.Fsused,
		Group:        blockdev.Group,
		Hctl:         blockdev.Hctl,
		Hotplug:      blockdev.Hotplug,
		Kname:        blockdev.Kname,
		Label:        blockdev.Label,
		LogSec:       blockdev.LogSec,
		MajMin:       blockdev.MajMin,
		MinIo:        blockdev.MinIo,
		Mode:         blockdev.Mode,
		Model:        blockdev.Model,
		Mountpoint:   blockdev.Mountpoint,
		OptIo:        blockdev.OptIo,
		Owner:        blockdev.Owner,
		Partflags:    blockdev.Partflags,
		Partlabel:    blockdev.Partlabel,
		Parttype:     blockdev.Parttype,
		Partuuid:     blockdev.Partuuid,
		Path:         blockdev.Path,
		PhySec:       blockdev.PhySec,
		Pttype:       blockdev.Pttype,
		Ptuuid:       blockdev.Ptuuid,
		Ra:           blockdev.Ra,
		Rand:         blockdev.Rand,
		Rev:          blockdev.Rev,
		Rm:           blockdev.Rm,
		Ro:           blockdev.Ro,
		Rota:         blockdev.Rota,
		RqSize:       blockdev.RqSize,
		Sched:        blockdev.Sched,
		Serial:       blockdev.Serial,
		Size:         blockdev.Size,
		State:        blockdev.State,
		Subsystems:   blockdev.Subsystems,
		Tran:         blockdev.Tran,
		Type:         blockdev.Type,
		UUID:         blockdev.UUID,
		Vendor:       blockdev.Vendor,
		Wsame:        blockdev.Wsame,
		WWN:          WWN,
	}

	metadata = graph.Metadata{
		"Index":      blockdev.getID(),
		"Path":       blockdev.getPath(),
		"MajorMinor": blockdev.MajMin,
		"Type":       nodeType,
		"Name":       blockdev.getName(),
		"Manager":    managerType,
		"Attributes": blockdevMetadata,
	}

	return metadata
}

func (p *ProbeHandler) findBlockDev(path string, id string, newDevInfo []BlockDevice) bool {
	for i := range newDevInfo {
		if id == newDevInfo[i].getID() && path == newDevInfo[i].getPath() {
			return true
		}
		if p.findBlockDev(path, id, newDevInfo[i].Children) == true {
			return true
		}
	}
	return false
}

func (p *ProbeHandler) getLVs() ([]LogicalVolume, error) {
	var (
		cmdOut []byte
		err    error
		report LVReport
	)

	lvsPath := p.Ctx.Config.GetString("agent.topology.blockdev.lvs_path")

	if lvsPath == "" {
		lvsPath = "/usr/sbin/lvs"
	}
	if cmdOut, err = exec.Command(lvsPath, "--reportformat", "json", "-o",
		"lv_name,lv_path,lv_dm_path,lv_layout,vg_name,lv_attr,lv_size,pool_lv,origin,data_percent,metadata_percent,move_pv,mirror_log,copy_percent,convert_lv,lv_kernel_major,lv_kernel_minor,lv_uuid").Output(); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(cmdOut, &report); err != nil {
		return nil, err
	}
	if len(report.Report) == 0 {
		return nil, nil
	}
	return report.Report[0].LVList, nil
}

func (p *ProbeHandler) getVGs() ([]VolumeGroup, error) {
	var (
		cmdOut []byte
		err    error
		report VGReport
	)

	vgsPath := p.Ctx.Config.GetString("agent.topology.blockdev.vgs_path")

	if vgsPath == "" {
		vgsPath = "/usr/sbin/vgs"
	}
	if cmdOut, err = exec.Command(vgsPath, "--reportformat", "json", "-o", "vg_name,vg_attr,vg_extent_size,pv_count,lv_count,snap_count,vg_size,vg_free,vg_uuid,vg_profile").Output(); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(cmdOut, &report); err != nil {
		return nil, err
	}
	if len(report.Report) == 0 {
		return nil, nil
	}
	return report.Report[0].LVList, nil
}

func (p *ProbeHandler) getPVs() ([]PhysicalVolume, error) {
	var (
		cmdOut []byte
		err    error
		report PVReport
	)

	pvsPath := p.Ctx.Config.GetString("agent.topology.blockdev.pvs_path")

	if pvsPath == "" {
		pvsPath = "/usr/sbin/pvs"
	}
	if cmdOut, err = exec.Command(pvsPath, "--reportformat", "json", "-o",
		"pv_name,vg_name,pv_fmt,pv_attr,pv_size,pv_free,dev_size,pv_uuid,pv_major,pv_minor").Output(); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(cmdOut, &report); err != nil {
		return nil, err
	}
	if len(report.Report) == 0 {
		return nil, nil
	}
	return report.Report[0].LVList, nil
}

func (p *ProbeHandler) getBlockDevices() ([]BlockDevice, error) {
	var (
		cmdOut []byte
		err    error
		intf   interface{}
		result Devices
	)

	lsblkPath := p.Ctx.Config.GetString("agent.topology.blockdev.lsblk_path")

	if lsblkPath == "" {
		lsblkPath = "/usr/bin/lsblk"
	}

	if cmdOut, err = exec.Command(lsblkPath, "-pO", "--json").Output(); err != nil {
		return nil, err
	}

	if err = json.Unmarshal([]byte(cmdOut[:]), &intf); err != nil {
		return nil, err
	}

	if err = mapstructure.WeakDecode(intf, &result); err != nil {
		return nil, err
	}

	return result.Blockdevices, nil
}

func (p *ProbeHandler) deleteIfRemoved(currentDevInfo blockdevInfo, newDevInfo []BlockDevice) {
	if p.findBlockDev(currentDevInfo.Path, currentDevInfo.ID, newDevInfo) == false {
		p.unregisterBlockdev(currentDevInfo.Path)
	}
}

func (p *ProbeHandler) registerLVMDevs(vgList []VolumeGroup, pvList []PhysicalVolume, lvList []LogicalVolume) *graph.Node {
	var lvMetaData graph.Metadata

	for _, vg := range vgList {
		_ = p.addGroupByName(vg.Name, vg.UUID)
		fmt.Println(vg)
	}

	for _, lv := range lvList {
		lvNode := p.Ctx.Graph.LookupFirstNode(graph.Metadata{"Path": lv.Path})

		if lvNode == nil {
			lvMetaData = p.getLVMetaData(lv)
			var err error
			if lvNode, err = p.Ctx.Graph.NewNode(graph.GenID(), lvMetaData); err != nil {
				p.Ctx.Logger.Error(err)
				return nil
			}
		}

		groupNode := p.Ctx.Graph.LookupFirstNode(graph.Metadata{"Name": lv.VG})

		topology.AddOwnershipLink(p.Ctx.Graph, groupNode, lvNode, nil)

		fmt.Println(lv)
	}
	for _, pv := range pvList {
		pvNode := p.Ctx.Graph.LookupFirstNode(graph.Metadata{"Name": pv.Name})

		if pvNode == nil {
			lvMetaData = p.getPVMetaData(pv)
			var err error
			if pvNode, err = p.Ctx.Graph.NewNode(graph.GenID(), lvMetaData); err != nil {
				p.Ctx.Logger.Error(err)
				return nil
			}
		}
		groupNode := p.Ctx.Graph.LookupFirstNode(graph.Metadata{"Name": pv.VG})

		topology.AddOwnershipLink(p.Ctx.Graph, groupNode, pvNode, nil)
		fmt.Println(pv)
	}

	return nil
}

func (p *ProbeHandler) registerBlockdev(blockdev BlockDevice, parentWWN string) []*graph.Node {
	var groupNodes []*graph.Node
	childCount := len(blockdev.Children)
	if childCount > 0 {
		for i := range blockdev.Children {
			groupNodes = p.registerBlockdev(blockdev.Children[i], blockdev.WWN)
		}
	} else {
		groupNodes = append(groupNodes, p.addGroup(blockdev, parentWWN))
	}

	p.Lock()
	defer p.Unlock()
	var graphMetaData graph.Metadata

	if _, ok := p.blockdevMap[blockdev.Name]; ok {
		return groupNodes
	}

	p.Ctx.Graph.Lock()
	defer p.Ctx.Graph.Unlock()

	node := p.Ctx.Graph.LookupFirstNode(graph.Metadata{"MajorMinor": blockdev.MajMin})

	if node == nil {
		graphMetaData = p.getMetaData(blockdev, childCount, parentWWN)
		var err error
		if node, err = p.Ctx.Graph.NewNode(graph.GenID(), graphMetaData); err != nil {
			p.Ctx.Logger.Error(err)
			return nil
		}

		// If there is a WWN - check to see if there are other paths to the same block device.
		if blockdev.WWN != "" {
			multiPathNodes := p.Ctx.Graph.GetNodes(graph.Metadata{"Index": blockdev.getID()})

			for _, multiPathNode := range multiPathNodes {
				topology.AddLink(p.Ctx.Graph, multiPathNode, node, "multipath", nil)
			}
		}

		for _, groupNode := range groupNodes {
			topology.AddLink(p.Ctx.Graph, groupNode, node, "ingroup", nil)
		}

		// Link physical disks and DVD/rom devices to the blockdevGroupName
		if blockdev.Type == "disk" || blockdev.Type == "rom" {
			linkMetadata := graph.Metadata{
				"Type": "blockdevlink",
			}
			topology.AddLink(p.Ctx.Graph, p.Groups[blockdevGroupName], node, "isablockdev", linkMetadata)
		}

	}

	for i := range blockdev.Children {
		childNode := p.Ctx.Graph.LookupFirstNode(graph.Metadata{"MajorMinor": blockdev.Children[i].MajMin})
		if childNode != nil {
			topology.AddLink(p.Ctx.Graph, childNode, node, "uses", nil)
		}
	}

	p.blockdevMap[blockdev.Name] = blockdevInfo{
		Name: blockdev.Name,
		ID:   blockdev.getID(),
		Path: blockdev.Path,
		Node: node,
	}
	return groupNodes
}

// Prior to version 2.33 lsblk generated JSON that encodes all values
// as strings.  Version 2.33 changed how integers and booleans are
// encoded.  unmarshalWeakTypeJSON() will decode either version of the
// JSON.
func (p *ProbeHandler) unmarshalWeakTypeJSON(jsonBytes []byte) ([]BlockDevice, error) {
	var decoder *mapstructure.Decoder
	var err error
	var result []BlockDevice

	deviceMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(jsonBytes), &deviceMap)
	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           &result,
	}
	if decoder, err = mapstructure.NewDecoder(config); err != nil {
		return nil, err
	}
	if err := decoder.Decode(deviceMap["blockdevices"]); err != nil {
		return nil, err
	}
	return result, nil
}

func (p *ProbeHandler) unregisterBlockdev(id string) {
	p.Lock()
	defer p.Unlock()

	info, ok := p.blockdevMap[id]
	if !ok {
		return
	}

	p.Ctx.Graph.Lock()
	if err := p.Ctx.Graph.DelNode(info.Node); err != nil {
		p.Ctx.Graph.Unlock()
		p.Ctx.Logger.Error(err)
		return
	}
	p.Ctx.Graph.Unlock()

	delete(p.blockdevMap, id)
}

func (p *ProbeHandler) connect() error {
	var (
		err    error
		result Devices
	)
	vgList, _ := p.getVGs()
	pvList, _ := p.getPVs()
	lvList, _ := p.getLVs()

	p.registerLVMDevs(vgList, pvList, lvList)

	blockDevices, err := p.getBlockDevices()

	if err != nil {
		return err
	}

	// loop through the devices in the current map to make sure they haven't
	// been removed
	for _, current := range p.blockdevMap {
		p.deleteIfRemoved(current, result.Blockdevices)
	}

	for _, blockdev := range blockDevices {
		p.registerBlockdev(blockdev, "")
	}

	return nil
}

func (p *ProbeHandler) newMetricsFromBlockdev(blockdevPath string) *BlockMetric {
	var (
		cmdOut []byte
		err    error
		stats  sysstat
		intf   interface{}
	)

	iostatPath := p.Ctx.Config.GetString("agent.topology.blockdev.iostat_path")

	if iostatPath == "" {
		iostatPath = "/usr/bin/iostat"
	}

	if cmdOut, err = exec.Command(iostatPath, "-dx", "-o", "JSON", blockdevPath).Output(); err != nil {
		return nil
	}

	if err = json.Unmarshal([]byte(cmdOut[:]), &intf); err != nil {
		return nil
	}

	if err = mapstructure.WeakDecode(intf, &stats); err != nil {
		return nil
	}

	if stats.Sysstat.Hosts == nil || len(stats.Sysstat.Hosts) == 0 || stats.Sysstat.Hosts[0].Statistics == nil ||
		len(stats.Sysstat.Hosts[0].Statistics) == 0 || stats.Sysstat.Hosts[0].Statistics[0].Metrics == nil ||
		len(stats.Sysstat.Hosts[0].Statistics[0].Metrics) == 0 {
		return nil
	}
	return stats.Sysstat.Hosts[0].Statistics[0].Metrics[0].MakeCopy()
}

func (p *ProbeHandler) addBlockDevData(now, last time.Time) {
	for _, blockdev := range p.blockdevMap {
		currMetric := p.newMetricsFromBlockdev(fmt.Sprintf("%v", blockdev.Path))
		if currMetric == nil {
			continue
		}
		currMetric.Last = int64(common.UnixMillis(now))
		p.Ctx.Graph.Lock()
		tr := p.Ctx.Graph.StartMetadataTransaction(blockdev.Node)
		tr.AddMetadata("BlockdevMetric", currMetric)
		tr.Commit()
		p.Ctx.Graph.Unlock()
	}
}

func (p *ProbeHandler) updateBlockDevMetric(now, last time.Time) {
	for _, blockdev := range p.blockdevMap {
		currMetric := p.newMetricsFromBlockdev(fmt.Sprintf("%v", blockdev.Path))
		if currMetric == nil {
			continue
		}
		currMetric.Last = int64(common.UnixMillis(now))
		p.Ctx.Graph.Lock()
		tr := p.Ctx.Graph.StartMetadataTransaction(blockdev.Node)
		tr.AddMetadata("BlockdevMetric", currMetric)
		tr.Commit()
		p.Ctx.Graph.Unlock()
	}
}

// Do adds a group for block devices, then issues a lsblk command and parsers the
// JSON output as a basis for the blockdev links.
func (p *ProbeHandler) Do(ctx context.Context, wg *sync.WaitGroup) error {
	p.addGroupByName(blockdevGroupName, "")
	p.Ctx.Graph.Lock()
	defer p.Ctx.Graph.Unlock()

	if !topology.HaveLink(p.Ctx.Graph, p.Ctx.RootNode, p.Groups[blockdevGroupName], "connected") {
		if _, err := topology.AddLink(p.Ctx.Graph, p.Ctx.RootNode, p.Groups[blockdevGroupName], "connected", nil); err != nil {
			p.Ctx.Logger.Error(err)
			return err
		}
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		metricTicker := time.NewTicker(time.Duration(30) * time.Second)
		var sampleInterval int
		if sampleInterval = p.Ctx.Config.GetInt("agent.topology.blockdev.iostat_update"); sampleInterval == 0 {
			sampleInterval = defaultIostatInterval
		}
		after := time.After(time.Duration(sampleInterval) * time.Second)
		last := time.Now().UTC()
		for {
			if err := p.connect(); err != nil {
				p.Ctx.Logger.Error(err)
				return
			}

			select {
			case t := <-metricTicker.C:
				now := t.UTC()
				p.updateBlockDevMetric(now, last)
				last = now
			case <-ctx.Done():
				return
			case <-after:
			}
		}
	}()
	return nil
}

// NewProbe initializes a new topology blockdev probe
func NewProbe(ctx tp.Context, bundle *probe.Bundle) (probe.Handler, error) {

	p := &ProbeHandler{
		blockdevMap: make(map[string]blockdevInfo),
		Groups:      make(map[string]*graph.Node),
		Ctx:         ctx,
	}

	return probes.NewProbeWrapper(p), nil
}

// Register registers graph metadata decoders
func Register() {
	graph.NodeMetadataDecoders["BlockdevMetric"] = MetricDecoder
}

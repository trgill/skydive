//go:generate go run ../../../scripts/gendecoder.go -package github.com/skydive-project/skydive/topology/probes/blockdev

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
	"encoding/json"
	"fmt"

	"github.com/skydive-project/skydive/common"
)

// Metadata describe the metadata of a docker container
// easyjson:json
// gendecoder
type Metadata struct {
	key                                   string
	storageID                             string
	storageLabel                          string
	partLabel                             string
	partUUID                              string
	storagePath                           string
	storageUUID                           string
	DEVTYPE                               string
	PWD                                   string
	SHLVL                                 string
	SUBSYSTEM                             string
	TAGS                                  string
	UDEVLOG_QUEUE_DISCARD_GRANULARIT      string
	UDEVLOG_QUEUE_DISCARD_MAX_BYTES       string
	UDEVLOG_QUEUE_DISCARD_ZEROES_DATA     string
	UDEVLOG_QUEUE_LOGICAL_BLOCK_SIZE      string
	UDEVLOG_QUEUE_MINIMUM_IO_SIZE         string
	UDEVLOG_QUEUE_OPTIMAL_IO_SIZE         string
	UDEVLOG_QUEUE_PHYSICAL_BLOCK_SIZE     string
	UDEVLOG_QUEUE_READ_AHEAD_KB           string
	UDEVLOG_QUEUE_ROTATIONAL              string
	UDEVLOG_BLOCK_ALIGNMENT_OFFSET        string
	UDEVLOG_BLOCK_DISCARD_ALIGNMENT       string
	UDEVLOG_BLOCK_REMOVABLE               string
	UDEVLOG_BLOCK_RO                      string
	UDEVLOG_LSBLK_MOUNTPOINT              string
	UDEVLOG_LSBLK_LABEL                   string
	UDEVLOG_LSBLK_STATE                   string
	UDEVLOG_LSBLK_OWNER                   string
	UDEVLOG_LSBLK_GROUP                   string
	UDEVLOG_LSBLK_MODE                    string
	UDEVLOG_BLOCK_EXT_RANGE               string
	_COMM                                 string
	_SELINUX_CONTEXT                      string
	_SYSTEMD_CGROUP                       string
	_SYSTEMD_UNIT                         string
	_SYSTEMD_INVOCATION_ID                string
	DM_ACTIVATION                         string
	DM_LV_LAYER                           string
	DM_NOSCAN                             string
	DM_SUBSYSTEM_UDEV_FLAG0               string
	DM_SUSPENDED                          string
	DM_UDEV_DISABLE_LIBRARY_FALLBACK_FLAG string
	DM_UDEV_DISABLE_OTHER_RULES_FLAG      string
	DM_UDEV_PRIMARY_SOURCE_FLAG           string
	DM_UDEV_RULES_VSN                     string
	DM_VG_NAME                            string
	MAJOR                                 string
	UDEVLOG_QUEUE_NR_REQUESTS             string
	UDEVLOG_QUEUE_SCHEDULER               string
	UDEVLOG_BLOCK_CAPABILITY              string
	UDEVLOG_LSBLK_TYPE                    string
	UDEVLOG_DM_READ_AHEAD                 string
	UDEVLOG_DM_ATTR                       string
	UDEVLOG_DM_TABLES_LOADED              string
	UDEVLOG_DM_READONLY                   string
	UDEVLOG_DM_SEGMENTS                   string
	UDEVLOG_DM_EVENTS                     string
	UDEVLOG_DM_SUBSYSTEM                  string
	UDEVLOG_DM_LV_LAYER                   string
	UDEVLOG_DM_TABLE_INACTIVE             string
	UDEVLOG_BLOCK_SIZE                    string
	DEVLINKS                              string
	DEVNAME                               string
	DEVPATH                               string
	DM_COOKIE                             string
	DM_DISABLE_OTHER_RULES_FLAG_OLD       string
	DM_LV_NAME                            string
	DM_NAME                               string
	DM_UUID                               string
	MINOR                                 string
	SEQNUM                                string
	USEC_INITIALIZED                      string
	PERSISTENT_STORAGE_ID                 string
	UDEVLOG_DM_OPEN                       string
	UDEVLOG_DM_DEVICE_COUNT               string
	UDEVLOG_DM_DEVS_USED                  string
	UDEVLOG_DM_DEVNOS_USED                string
	UDEVLOG_DM_BLKDEVS_USED               string
	UDEVLOG_DM_DEVICE_REF_COUNT           string
	UDEVLOG_DM_NAMES_USING_DEV            string
	UDEVLOG_DM_DEVNOS_USING_DEV           string
	UDEVLOG_DM_TABLE_LIVE                 string
}

// MetadataDecoder implements a json message raw decoder
func MetadataDecoder(raw json.RawMessage) (common.Getter, error) {
	var m Metadata
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, fmt.Errorf("unable to unmarshal docker metadata %s: %s", string(raw), err)
	}

	return &m, nil
}

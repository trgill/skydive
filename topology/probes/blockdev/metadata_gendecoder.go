// Code generated - DO NOT EDIT.

package blockdev

import (
	"github.com/skydive-project/skydive/common"
	"strings"
)

func (obj *Metadata) GetFieldBool(key string) (bool, error) {
	return false, common.ErrFieldNotFound
}

func (obj *Metadata) GetFieldInt64(key string) (int64, error) {
	return 0, common.ErrFieldNotFound
}

func (obj *Metadata) GetFieldString(key string) (string, error) {
	switch key {
	case "key":
		return string(obj.key), nil
	case "storageID":
		return string(obj.storageID), nil
	case "storageLabel":
		return string(obj.storageLabel), nil
	case "partLabel":
		return string(obj.partLabel), nil
	case "partUUID":
		return string(obj.partUUID), nil
	case "storagePath":
		return string(obj.storagePath), nil
	case "storageUUID":
		return string(obj.storageUUID), nil
	case "DEVTYPE":
		return string(obj.DEVTYPE), nil
	case "PWD":
		return string(obj.PWD), nil
	case "SHLVL":
		return string(obj.SHLVL), nil
	case "SUBSYSTEM":
		return string(obj.SUBSYSTEM), nil
	case "TAGS":
		return string(obj.TAGS), nil
	case "UDEVLOG_QUEUE_DISCARD_GRANULARIT":
		return string(obj.UDEVLOG_QUEUE_DISCARD_GRANULARIT), nil
	case "UDEVLOG_QUEUE_DISCARD_MAX_BYTES":
		return string(obj.UDEVLOG_QUEUE_DISCARD_MAX_BYTES), nil
	case "UDEVLOG_QUEUE_DISCARD_ZEROES_DATA":
		return string(obj.UDEVLOG_QUEUE_DISCARD_ZEROES_DATA), nil
	case "UDEVLOG_QUEUE_LOGICAL_BLOCK_SIZE":
		return string(obj.UDEVLOG_QUEUE_LOGICAL_BLOCK_SIZE), nil
	case "UDEVLOG_QUEUE_MINIMUM_IO_SIZE":
		return string(obj.UDEVLOG_QUEUE_MINIMUM_IO_SIZE), nil
	case "UDEVLOG_QUEUE_OPTIMAL_IO_SIZE":
		return string(obj.UDEVLOG_QUEUE_OPTIMAL_IO_SIZE), nil
	case "UDEVLOG_QUEUE_PHYSICAL_BLOCK_SIZE":
		return string(obj.UDEVLOG_QUEUE_PHYSICAL_BLOCK_SIZE), nil
	case "UDEVLOG_QUEUE_READ_AHEAD_KB":
		return string(obj.UDEVLOG_QUEUE_READ_AHEAD_KB), nil
	case "UDEVLOG_QUEUE_ROTATIONAL":
		return string(obj.UDEVLOG_QUEUE_ROTATIONAL), nil
	case "UDEVLOG_BLOCK_ALIGNMENT_OFFSET":
		return string(obj.UDEVLOG_BLOCK_ALIGNMENT_OFFSET), nil
	case "UDEVLOG_BLOCK_DISCARD_ALIGNMENT":
		return string(obj.UDEVLOG_BLOCK_DISCARD_ALIGNMENT), nil
	case "UDEVLOG_BLOCK_REMOVABLE":
		return string(obj.UDEVLOG_BLOCK_REMOVABLE), nil
	case "UDEVLOG_BLOCK_RO":
		return string(obj.UDEVLOG_BLOCK_RO), nil
	case "UDEVLOG_LSBLK_MOUNTPOINT":
		return string(obj.UDEVLOG_LSBLK_MOUNTPOINT), nil
	case "UDEVLOG_LSBLK_LABEL":
		return string(obj.UDEVLOG_LSBLK_LABEL), nil
	case "UDEVLOG_LSBLK_STATE":
		return string(obj.UDEVLOG_LSBLK_STATE), nil
	case "UDEVLOG_LSBLK_OWNER":
		return string(obj.UDEVLOG_LSBLK_OWNER), nil
	case "UDEVLOG_LSBLK_GROUP":
		return string(obj.UDEVLOG_LSBLK_GROUP), nil
	case "UDEVLOG_LSBLK_MODE":
		return string(obj.UDEVLOG_LSBLK_MODE), nil
	case "UDEVLOG_BLOCK_EXT_RANGE":
		return string(obj.UDEVLOG_BLOCK_EXT_RANGE), nil
	case "_COMM":
		return string(obj._COMM), nil
	case "_SELINUX_CONTEXT":
		return string(obj._SELINUX_CONTEXT), nil
	case "_SYSTEMD_CGROUP":
		return string(obj._SYSTEMD_CGROUP), nil
	case "_SYSTEMD_UNIT":
		return string(obj._SYSTEMD_UNIT), nil
	case "_SYSTEMD_INVOCATION_ID":
		return string(obj._SYSTEMD_INVOCATION_ID), nil
	case "DM_ACTIVATION":
		return string(obj.DM_ACTIVATION), nil
	case "DM_LV_LAYER":
		return string(obj.DM_LV_LAYER), nil
	case "DM_NOSCAN":
		return string(obj.DM_NOSCAN), nil
	case "DM_SUBSYSTEM_UDEV_FLAG0":
		return string(obj.DM_SUBSYSTEM_UDEV_FLAG0), nil
	case "DM_SUSPENDED":
		return string(obj.DM_SUSPENDED), nil
	case "DM_UDEV_DISABLE_LIBRARY_FALLBACK_FLAG":
		return string(obj.DM_UDEV_DISABLE_LIBRARY_FALLBACK_FLAG), nil
	case "DM_UDEV_DISABLE_OTHER_RULES_FLAG":
		return string(obj.DM_UDEV_DISABLE_OTHER_RULES_FLAG), nil
	case "DM_UDEV_PRIMARY_SOURCE_FLAG":
		return string(obj.DM_UDEV_PRIMARY_SOURCE_FLAG), nil
	case "DM_UDEV_RULES_VSN":
		return string(obj.DM_UDEV_RULES_VSN), nil
	case "DM_VG_NAME":
		return string(obj.DM_VG_NAME), nil
	case "MAJOR":
		return string(obj.MAJOR), nil
	case "UDEVLOG_QUEUE_NR_REQUESTS":
		return string(obj.UDEVLOG_QUEUE_NR_REQUESTS), nil
	case "UDEVLOG_QUEUE_SCHEDULER":
		return string(obj.UDEVLOG_QUEUE_SCHEDULER), nil
	case "UDEVLOG_BLOCK_CAPABILITY":
		return string(obj.UDEVLOG_BLOCK_CAPABILITY), nil
	case "UDEVLOG_LSBLK_TYPE":
		return string(obj.UDEVLOG_LSBLK_TYPE), nil
	case "UDEVLOG_DM_READ_AHEAD":
		return string(obj.UDEVLOG_DM_READ_AHEAD), nil
	case "UDEVLOG_DM_ATTR":
		return string(obj.UDEVLOG_DM_ATTR), nil
	case "UDEVLOG_DM_TABLES_LOADED":
		return string(obj.UDEVLOG_DM_TABLES_LOADED), nil
	case "UDEVLOG_DM_READONLY":
		return string(obj.UDEVLOG_DM_READONLY), nil
	case "UDEVLOG_DM_SEGMENTS":
		return string(obj.UDEVLOG_DM_SEGMENTS), nil
	case "UDEVLOG_DM_EVENTS":
		return string(obj.UDEVLOG_DM_EVENTS), nil
	case "UDEVLOG_DM_SUBSYSTEM":
		return string(obj.UDEVLOG_DM_SUBSYSTEM), nil
	case "UDEVLOG_DM_LV_LAYER":
		return string(obj.UDEVLOG_DM_LV_LAYER), nil
	case "UDEVLOG_DM_TABLE_INACTIVE":
		return string(obj.UDEVLOG_DM_TABLE_INACTIVE), nil
	case "UDEVLOG_BLOCK_SIZE":
		return string(obj.UDEVLOG_BLOCK_SIZE), nil
	case "DEVLINKS":
		return string(obj.DEVLINKS), nil
	case "DEVNAME":
		return string(obj.DEVNAME), nil
	case "DEVPATH":
		return string(obj.DEVPATH), nil
	case "DM_COOKIE":
		return string(obj.DM_COOKIE), nil
	case "DM_DISABLE_OTHER_RULES_FLAG_OLD":
		return string(obj.DM_DISABLE_OTHER_RULES_FLAG_OLD), nil
	case "DM_LV_NAME":
		return string(obj.DM_LV_NAME), nil
	case "DM_NAME":
		return string(obj.DM_NAME), nil
	case "DM_UUID":
		return string(obj.DM_UUID), nil
	case "MINOR":
		return string(obj.MINOR), nil
	case "SEQNUM":
		return string(obj.SEQNUM), nil
	case "USEC_INITIALIZED":
		return string(obj.USEC_INITIALIZED), nil
	case "PERSISTENT_STORAGE_ID":
		return string(obj.PERSISTENT_STORAGE_ID), nil
	case "UDEVLOG_DM_OPEN":
		return string(obj.UDEVLOG_DM_OPEN), nil
	case "UDEVLOG_DM_DEVICE_COUNT":
		return string(obj.UDEVLOG_DM_DEVICE_COUNT), nil
	case "UDEVLOG_DM_DEVS_USED":
		return string(obj.UDEVLOG_DM_DEVS_USED), nil
	case "UDEVLOG_DM_DEVNOS_USED":
		return string(obj.UDEVLOG_DM_DEVNOS_USED), nil
	case "UDEVLOG_DM_BLKDEVS_USED":
		return string(obj.UDEVLOG_DM_BLKDEVS_USED), nil
	case "UDEVLOG_DM_DEVICE_REF_COUNT":
		return string(obj.UDEVLOG_DM_DEVICE_REF_COUNT), nil
	case "UDEVLOG_DM_NAMES_USING_DEV":
		return string(obj.UDEVLOG_DM_NAMES_USING_DEV), nil
	case "UDEVLOG_DM_DEVNOS_USING_DEV":
		return string(obj.UDEVLOG_DM_DEVNOS_USING_DEV), nil
	case "UDEVLOG_DM_TABLE_LIVE":
		return string(obj.UDEVLOG_DM_TABLE_LIVE), nil
	}
	return "", common.ErrFieldNotFound
}

func (obj *Metadata) GetFieldKeys() []string {
	return []string{
		"key",
		"storageID",
		"storageLabel",
		"partLabel",
		"partUUID",
		"storagePath",
		"storageUUID",
		"DEVTYPE",
		"PWD",
		"SHLVL",
		"SUBSYSTEM",
		"TAGS",
		"UDEVLOG_QUEUE_DISCARD_GRANULARIT",
		"UDEVLOG_QUEUE_DISCARD_MAX_BYTES",
		"UDEVLOG_QUEUE_DISCARD_ZEROES_DATA",
		"UDEVLOG_QUEUE_LOGICAL_BLOCK_SIZE",
		"UDEVLOG_QUEUE_MINIMUM_IO_SIZE",
		"UDEVLOG_QUEUE_OPTIMAL_IO_SIZE",
		"UDEVLOG_QUEUE_PHYSICAL_BLOCK_SIZE",
		"UDEVLOG_QUEUE_READ_AHEAD_KB",
		"UDEVLOG_QUEUE_ROTATIONAL",
		"UDEVLOG_BLOCK_ALIGNMENT_OFFSET",
		"UDEVLOG_BLOCK_DISCARD_ALIGNMENT",
		"UDEVLOG_BLOCK_REMOVABLE",
		"UDEVLOG_BLOCK_RO",
		"UDEVLOG_LSBLK_MOUNTPOINT",
		"UDEVLOG_LSBLK_LABEL",
		"UDEVLOG_LSBLK_STATE",
		"UDEVLOG_LSBLK_OWNER",
		"UDEVLOG_LSBLK_GROUP",
		"UDEVLOG_LSBLK_MODE",
		"UDEVLOG_BLOCK_EXT_RANGE",
		"_COMM",
		"_SELINUX_CONTEXT",
		"_SYSTEMD_CGROUP",
		"_SYSTEMD_UNIT",
		"_SYSTEMD_INVOCATION_ID",
		"DM_ACTIVATION",
		"DM_LV_LAYER",
		"DM_NOSCAN",
		"DM_SUBSYSTEM_UDEV_FLAG0",
		"DM_SUSPENDED",
		"DM_UDEV_DISABLE_LIBRARY_FALLBACK_FLAG",
		"DM_UDEV_DISABLE_OTHER_RULES_FLAG",
		"DM_UDEV_PRIMARY_SOURCE_FLAG",
		"DM_UDEV_RULES_VSN",
		"DM_VG_NAME",
		"MAJOR",
		"UDEVLOG_QUEUE_NR_REQUESTS",
		"UDEVLOG_QUEUE_SCHEDULER",
		"UDEVLOG_BLOCK_CAPABILITY",
		"UDEVLOG_LSBLK_TYPE",
		"UDEVLOG_DM_READ_AHEAD",
		"UDEVLOG_DM_ATTR",
		"UDEVLOG_DM_TABLES_LOADED",
		"UDEVLOG_DM_READONLY",
		"UDEVLOG_DM_SEGMENTS",
		"UDEVLOG_DM_EVENTS",
		"UDEVLOG_DM_SUBSYSTEM",
		"UDEVLOG_DM_LV_LAYER",
		"UDEVLOG_DM_TABLE_INACTIVE",
		"UDEVLOG_BLOCK_SIZE",
		"DEVLINKS",
		"DEVNAME",
		"DEVPATH",
		"DM_COOKIE",
		"DM_DISABLE_OTHER_RULES_FLAG_OLD",
		"DM_LV_NAME",
		"DM_NAME",
		"DM_UUID",
		"MINOR",
		"SEQNUM",
		"USEC_INITIALIZED",
		"PERSISTENT_STORAGE_ID",
		"UDEVLOG_DM_OPEN",
		"UDEVLOG_DM_DEVICE_COUNT",
		"UDEVLOG_DM_DEVS_USED",
		"UDEVLOG_DM_DEVNOS_USED",
		"UDEVLOG_DM_BLKDEVS_USED",
		"UDEVLOG_DM_DEVICE_REF_COUNT",
		"UDEVLOG_DM_NAMES_USING_DEV",
		"UDEVLOG_DM_DEVNOS_USING_DEV",
		"UDEVLOG_DM_TABLE_LIVE",
	}
}

func (obj *Metadata) MatchBool(key string, predicate common.BoolPredicate) bool {
	return false
}

func (obj *Metadata) MatchInt64(key string, predicate common.Int64Predicate) bool {
	return false
}

func (obj *Metadata) MatchString(key string, predicate common.StringPredicate) bool {
	if b, err := obj.GetFieldString(key); err == nil {
		return predicate(b)
	}
	return false
}

func (obj *Metadata) GetField(key string) (interface{}, error) {
	if s, err := obj.GetFieldString(key); err == nil {
		return s, nil
	}
	return nil, common.ErrFieldNotFound
}

func init() {
	strings.Index("", ".")
}

// Code generated - DO NOT EDIT.

package blockdev

import (
	"github.com/skydive-project/skydive/common"
	"strings"
)

func (obj *Metadata) GetFieldBool(key string) (bool, error) {
	switch key {
	case "IsInitialized":
		return obj.IsInitialized, nil
	}

	return false, common.ErrFieldNotFound
}

func (obj *Metadata) GetFieldInt64(key string) (int64, error) {
	return 0, common.ErrFieldNotFound
}

func (obj *Metadata) GetFieldString(key string) (string, error) {
	switch key {
	case "BlockDevID":
		return string(obj.BlockDevID), nil
	case "BlockDevName":
		return string(obj.BlockDevName), nil
	case "Sysname":
		return string(obj.Sysname), nil
	case "Syspath":
		return string(obj.Syspath), nil
	case "Devpath":
		return string(obj.Devpath), nil
	case "Devnode":
		return string(obj.Devnode), nil
	case "Subsystem":
		return string(obj.Subsystem), nil
	case "Devtype":
		return string(obj.Devtype), nil
	case "Sysnum":
		return string(obj.Sysnum), nil
	case "Driver":
		return string(obj.Driver), nil
	case "DevLinks":
		return string(obj.DevLinks), nil
	}
	return "", common.ErrFieldNotFound
}

func (obj *Metadata) GetFieldKeys() []string {
	return []string{
		"BlockDevID",
		"BlockDevName",
		"Sysname",
		"Syspath",
		"Devpath",
		"Devnode",
		"Subsystem",
		"Devtype",
		"Sysnum",
		"IsInitialized",
		"Driver",
		"DevLinks",
		"Labels",
	}
}

func (obj *Metadata) MatchBool(key string, predicate common.BoolPredicate) bool {
	if b, err := obj.GetFieldBool(key); err == nil {
		return predicate(b)
	}

	first := key
	index := strings.Index(key, ".")
	if index != -1 {
		first = key[:index]
	}

	switch first {
	case "Labels":
		if index != -1 && obj.Labels != nil {
			return obj.Labels.MatchBool(key[index+1:], predicate)
		}
	}
	return false
}

func (obj *Metadata) MatchInt64(key string, predicate common.Int64Predicate) bool {
	first := key
	index := strings.Index(key, ".")
	if index != -1 {
		first = key[:index]
	}

	switch first {

	case "Labels":
		if index != -1 && obj.Labels != nil {
			return obj.Labels.MatchInt64(key[index+1:], predicate)
		}
	}
	return false
}

func (obj *Metadata) MatchString(key string, predicate common.StringPredicate) bool {
	if b, err := obj.GetFieldString(key); err == nil {
		return predicate(b)
	}

	first := key
	index := strings.Index(key, ".")
	if index != -1 {
		first = key[:index]
	}

	switch first {

	case "Labels":
		if index != -1 && obj.Labels != nil {
			return obj.Labels.MatchString(key[index+1:], predicate)
		}
	}
	return false
}

func (obj *Metadata) GetField(key string) (interface{}, error) {
	if s, err := obj.GetFieldString(key); err == nil {
		return s, nil
	}

	if b, err := obj.GetFieldBool(key); err == nil {
		return b, nil
	}

	first := key
	index := strings.Index(key, ".")
	if index != -1 {
		first = key[:index]
	}

	switch first {
	case "Labels":
		if obj.Labels != nil {
			if index != -1 {
				return obj.Labels.GetField(key[index+1:])
			} else {
				return obj.Labels, nil
			}
		}

	}
	return nil, common.ErrFieldNotFound
}

func init() {
	strings.Index("", ".")
}

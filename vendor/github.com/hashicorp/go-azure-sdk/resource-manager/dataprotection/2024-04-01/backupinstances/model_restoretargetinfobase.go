package backupinstances

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RestoreTargetInfoBase interface {
	RestoreTargetInfoBase() BaseRestoreTargetInfoBaseImpl
}

var _ RestoreTargetInfoBase = BaseRestoreTargetInfoBaseImpl{}

type BaseRestoreTargetInfoBaseImpl struct {
	ObjectType      string         `json:"objectType"`
	RecoveryOption  RecoveryOption `json:"recoveryOption"`
	RestoreLocation *string        `json:"restoreLocation,omitempty"`
}

func (s BaseRestoreTargetInfoBaseImpl) RestoreTargetInfoBase() BaseRestoreTargetInfoBaseImpl {
	return s
}

var _ RestoreTargetInfoBase = RawRestoreTargetInfoBaseImpl{}

// RawRestoreTargetInfoBaseImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawRestoreTargetInfoBaseImpl struct {
	restoreTargetInfoBase BaseRestoreTargetInfoBaseImpl
	Type                  string
	Values                map[string]interface{}
}

func (s RawRestoreTargetInfoBaseImpl) RestoreTargetInfoBase() BaseRestoreTargetInfoBaseImpl {
	return s.restoreTargetInfoBase
}

func UnmarshalRestoreTargetInfoBaseImplementation(input []byte) (RestoreTargetInfoBase, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling RestoreTargetInfoBase into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["objectType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "ItemLevelRestoreTargetInfo") {
		var out ItemLevelRestoreTargetInfo
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ItemLevelRestoreTargetInfo: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "RestoreFilesTargetInfo") {
		var out RestoreFilesTargetInfo
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into RestoreFilesTargetInfo: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "RestoreTargetInfo") {
		var out RestoreTargetInfo
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into RestoreTargetInfo: %+v", err)
		}
		return out, nil
	}

	var parent BaseRestoreTargetInfoBaseImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseRestoreTargetInfoBaseImpl: %+v", err)
	}

	return RawRestoreTargetInfoBaseImpl{
		restoreTargetInfoBase: parent,
		Type:                  value,
		Values:                temp,
	}, nil

}

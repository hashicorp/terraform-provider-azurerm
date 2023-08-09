package backupinstances

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RestoreTargetInfoBase interface {
}

// RawModeOfTransitImpl is returned when the Discriminated Value
// doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawRestoreTargetInfoBaseImpl struct {
	Type   string
	Values map[string]interface{}
}

func unmarshalRestoreTargetInfoBaseImplementation(input []byte) (RestoreTargetInfoBase, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling RestoreTargetInfoBase into map[string]interface: %+v", err)
	}

	value, ok := temp["objectType"].(string)
	if !ok {
		return nil, nil
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

	out := RawRestoreTargetInfoBaseImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}

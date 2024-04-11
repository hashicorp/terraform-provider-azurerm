package appplatform

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AcceleratorAuthSetting interface {
}

// RawAcceleratorAuthSettingImpl is returned when the Discriminated Value
// doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawAcceleratorAuthSettingImpl struct {
	Type   string
	Values map[string]interface{}
}

func unmarshalAcceleratorAuthSettingImplementation(input []byte) (AcceleratorAuthSetting, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling AcceleratorAuthSetting into map[string]interface: %+v", err)
	}

	value, ok := temp["authType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "BasicAuth") {
		var out AcceleratorBasicAuthSetting
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AcceleratorBasicAuthSetting: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Public") {
		var out AcceleratorPublicSetting
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AcceleratorPublicSetting: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SSH") {
		var out AcceleratorSshSetting
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AcceleratorSshSetting: %+v", err)
		}
		return out, nil
	}

	out := RawAcceleratorAuthSettingImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}

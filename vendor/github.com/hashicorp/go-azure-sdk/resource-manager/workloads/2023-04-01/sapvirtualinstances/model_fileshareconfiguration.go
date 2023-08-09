package sapvirtualinstances

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FileShareConfiguration interface {
}

// RawModeOfTransitImpl is returned when the Discriminated Value
// doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawFileShareConfigurationImpl struct {
	Type   string
	Values map[string]interface{}
}

func unmarshalFileShareConfigurationImplementation(input []byte) (FileShareConfiguration, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling FileShareConfiguration into map[string]interface: %+v", err)
	}

	value, ok := temp["configurationType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "CreateAndMount") {
		var out CreateAndMountFileShareConfiguration
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CreateAndMountFileShareConfiguration: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Mount") {
		var out MountFileShareConfiguration
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MountFileShareConfiguration: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Skip") {
		var out SkipFileShareConfiguration
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SkipFileShareConfiguration: %+v", err)
		}
		return out, nil
	}

	out := RawFileShareConfigurationImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}

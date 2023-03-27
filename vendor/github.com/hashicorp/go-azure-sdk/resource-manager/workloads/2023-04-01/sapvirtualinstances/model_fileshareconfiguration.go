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

	type RawFileShareConfigurationImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawFileShareConfigurationImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}

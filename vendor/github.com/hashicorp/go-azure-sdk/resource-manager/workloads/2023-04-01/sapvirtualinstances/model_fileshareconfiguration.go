package sapvirtualinstances

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FileShareConfiguration interface {
	FileShareConfiguration() BaseFileShareConfigurationImpl
}

var _ FileShareConfiguration = BaseFileShareConfigurationImpl{}

type BaseFileShareConfigurationImpl struct {
	ConfigurationType ConfigurationType `json:"configurationType"`
}

func (s BaseFileShareConfigurationImpl) FileShareConfiguration() BaseFileShareConfigurationImpl {
	return s
}

var _ FileShareConfiguration = RawFileShareConfigurationImpl{}

// RawFileShareConfigurationImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawFileShareConfigurationImpl struct {
	fileShareConfiguration BaseFileShareConfigurationImpl
	Type                   string
	Values                 map[string]interface{}
}

func (s RawFileShareConfigurationImpl) FileShareConfiguration() BaseFileShareConfigurationImpl {
	return s.fileShareConfiguration
}

func UnmarshalFileShareConfigurationImplementation(input []byte) (FileShareConfiguration, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling FileShareConfiguration into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["configurationType"]; ok {
		value = fmt.Sprintf("%v", v)
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

	var parent BaseFileShareConfigurationImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseFileShareConfigurationImpl: %+v", err)
	}

	return RawFileShareConfigurationImpl{
		fileShareConfiguration: parent,
		Type:                   value,
		Values:                 temp,
	}, nil

}

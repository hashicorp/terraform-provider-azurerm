package sapvirtualinstances

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InfrastructureConfiguration interface {
}

func unmarshalInfrastructureConfigurationImplementation(input []byte) (InfrastructureConfiguration, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling InfrastructureConfiguration into map[string]interface: %+v", err)
	}

	value, ok := temp["deploymentType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "SingleServer") {
		var out SingleServerConfiguration
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SingleServerConfiguration: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ThreeTier") {
		var out ThreeTierConfiguration
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ThreeTierConfiguration: %+v", err)
		}
		return out, nil
	}

	type RawInfrastructureConfigurationImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawInfrastructureConfigurationImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}

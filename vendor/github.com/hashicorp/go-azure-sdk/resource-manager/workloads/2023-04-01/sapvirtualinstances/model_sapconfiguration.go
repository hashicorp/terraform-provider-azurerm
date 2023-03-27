package sapvirtualinstances

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SAPConfiguration interface {
}

func unmarshalSAPConfigurationImplementation(input []byte) (SAPConfiguration, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling SAPConfiguration into map[string]interface: %+v", err)
	}

	value, ok := temp["configurationType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "Deployment") {
		var out DeploymentConfiguration
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DeploymentConfiguration: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DeploymentWithOSConfig") {
		var out DeploymentWithOSConfiguration
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DeploymentWithOSConfiguration: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Discovery") {
		var out DiscoveryConfiguration
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DiscoveryConfiguration: %+v", err)
		}
		return out, nil
	}

	type RawSAPConfigurationImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawSAPConfigurationImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}

package replicationnetworkmappings

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkMappingFabricSpecificSettings interface {
}

func unmarshalNetworkMappingFabricSpecificSettingsImplementation(input []byte) (NetworkMappingFabricSpecificSettings, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling NetworkMappingFabricSpecificSettings into map[string]interface: %+v", err)
	}

	value, ok := temp["instanceType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "AzureToAzure") {
		var out AzureToAzureNetworkMappingSettings
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureToAzureNetworkMappingSettings: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "VmmToAzure") {
		var out VmmToAzureNetworkMappingSettings
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into VmmToAzureNetworkMappingSettings: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "VmmToVmm") {
		var out VmmToVmmNetworkMappingSettings
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into VmmToVmmNetworkMappingSettings: %+v", err)
		}
		return out, nil
	}

	type RawNetworkMappingFabricSpecificSettingsImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawNetworkMappingFabricSpecificSettingsImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}

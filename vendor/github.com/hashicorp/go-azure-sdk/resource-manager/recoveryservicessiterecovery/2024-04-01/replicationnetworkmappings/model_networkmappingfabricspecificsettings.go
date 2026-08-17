package replicationnetworkmappings

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkMappingFabricSpecificSettings interface {
	NetworkMappingFabricSpecificSettings() BaseNetworkMappingFabricSpecificSettingsImpl
}

var _ NetworkMappingFabricSpecificSettings = BaseNetworkMappingFabricSpecificSettingsImpl{}

type BaseNetworkMappingFabricSpecificSettingsImpl struct {
	InstanceType string `json:"instanceType"`
}

func (s BaseNetworkMappingFabricSpecificSettingsImpl) NetworkMappingFabricSpecificSettings() BaseNetworkMappingFabricSpecificSettingsImpl {
	return s
}

var _ NetworkMappingFabricSpecificSettings = RawNetworkMappingFabricSpecificSettingsImpl{}

// RawNetworkMappingFabricSpecificSettingsImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawNetworkMappingFabricSpecificSettingsImpl struct {
	networkMappingFabricSpecificSettings BaseNetworkMappingFabricSpecificSettingsImpl
	Type                                 string
	Values                               map[string]interface{}
}

func (s RawNetworkMappingFabricSpecificSettingsImpl) NetworkMappingFabricSpecificSettings() BaseNetworkMappingFabricSpecificSettingsImpl {
	return s.networkMappingFabricSpecificSettings
}

func (s RawNetworkMappingFabricSpecificSettingsImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalNetworkMappingFabricSpecificSettingsImplementation(input []byte) (NetworkMappingFabricSpecificSettings, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling NetworkMappingFabricSpecificSettings into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["instanceType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "AzureToAzure") {
		var out AzureToAzureNetworkMappingSettings
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureToAzureNetworkMappingSettings: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "VmmToAzure") {
		var out VMmToAzureNetworkMappingSettings
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into VMmToAzureNetworkMappingSettings: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "VmmToVmm") {
		var out VMmToVMmNetworkMappingSettings
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into VMmToVMmNetworkMappingSettings: %+v", err)
		}
		return out, nil
	}

	var parent BaseNetworkMappingFabricSpecificSettingsImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseNetworkMappingFabricSpecificSettingsImpl: %+v", err)
	}

	return RawNetworkMappingFabricSpecificSettingsImpl{
		networkMappingFabricSpecificSettings: parent,
		Type:                                 value,
		Values:                               temp,
	}, nil

}

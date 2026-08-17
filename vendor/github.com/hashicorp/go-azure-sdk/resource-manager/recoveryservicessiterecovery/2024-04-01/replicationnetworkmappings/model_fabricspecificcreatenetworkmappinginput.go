package replicationnetworkmappings

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FabricSpecificCreateNetworkMappingInput interface {
	FabricSpecificCreateNetworkMappingInput() BaseFabricSpecificCreateNetworkMappingInputImpl
}

var _ FabricSpecificCreateNetworkMappingInput = BaseFabricSpecificCreateNetworkMappingInputImpl{}

type BaseFabricSpecificCreateNetworkMappingInputImpl struct {
	InstanceType string `json:"instanceType"`
}

func (s BaseFabricSpecificCreateNetworkMappingInputImpl) FabricSpecificCreateNetworkMappingInput() BaseFabricSpecificCreateNetworkMappingInputImpl {
	return s
}

var _ FabricSpecificCreateNetworkMappingInput = RawFabricSpecificCreateNetworkMappingInputImpl{}

// RawFabricSpecificCreateNetworkMappingInputImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawFabricSpecificCreateNetworkMappingInputImpl struct {
	fabricSpecificCreateNetworkMappingInput BaseFabricSpecificCreateNetworkMappingInputImpl
	Type                                    string
	Values                                  map[string]interface{}
}

func (s RawFabricSpecificCreateNetworkMappingInputImpl) FabricSpecificCreateNetworkMappingInput() BaseFabricSpecificCreateNetworkMappingInputImpl {
	return s.fabricSpecificCreateNetworkMappingInput
}

func (s RawFabricSpecificCreateNetworkMappingInputImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalFabricSpecificCreateNetworkMappingInputImplementation(input []byte) (FabricSpecificCreateNetworkMappingInput, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling FabricSpecificCreateNetworkMappingInput into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["instanceType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "AzureToAzure") {
		var out AzureToAzureCreateNetworkMappingInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureToAzureCreateNetworkMappingInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "VmmToAzure") {
		var out VMmToAzureCreateNetworkMappingInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into VMmToAzureCreateNetworkMappingInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "VmmToVmm") {
		var out VMmToVMmCreateNetworkMappingInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into VMmToVMmCreateNetworkMappingInput: %+v", err)
		}
		return out, nil
	}

	var parent BaseFabricSpecificCreateNetworkMappingInputImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseFabricSpecificCreateNetworkMappingInputImpl: %+v", err)
	}

	return RawFabricSpecificCreateNetworkMappingInputImpl{
		fabricSpecificCreateNetworkMappingInput: parent,
		Type:                                    value,
		Values:                                  temp,
	}, nil

}

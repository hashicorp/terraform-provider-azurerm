package replicationnetworkmappings

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FabricSpecificUpdateNetworkMappingInput interface {
	FabricSpecificUpdateNetworkMappingInput() BaseFabricSpecificUpdateNetworkMappingInputImpl
}

var _ FabricSpecificUpdateNetworkMappingInput = BaseFabricSpecificUpdateNetworkMappingInputImpl{}

type BaseFabricSpecificUpdateNetworkMappingInputImpl struct {
	InstanceType string `json:"instanceType"`
}

func (s BaseFabricSpecificUpdateNetworkMappingInputImpl) FabricSpecificUpdateNetworkMappingInput() BaseFabricSpecificUpdateNetworkMappingInputImpl {
	return s
}

var _ FabricSpecificUpdateNetworkMappingInput = RawFabricSpecificUpdateNetworkMappingInputImpl{}

// RawFabricSpecificUpdateNetworkMappingInputImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawFabricSpecificUpdateNetworkMappingInputImpl struct {
	fabricSpecificUpdateNetworkMappingInput BaseFabricSpecificUpdateNetworkMappingInputImpl
	Type                                    string
	Values                                  map[string]interface{}
}

func (s RawFabricSpecificUpdateNetworkMappingInputImpl) FabricSpecificUpdateNetworkMappingInput() BaseFabricSpecificUpdateNetworkMappingInputImpl {
	return s.fabricSpecificUpdateNetworkMappingInput
}

func UnmarshalFabricSpecificUpdateNetworkMappingInputImplementation(input []byte) (FabricSpecificUpdateNetworkMappingInput, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling FabricSpecificUpdateNetworkMappingInput into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["instanceType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "AzureToAzure") {
		var out AzureToAzureUpdateNetworkMappingInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureToAzureUpdateNetworkMappingInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "VmmToAzure") {
		var out VMmToAzureUpdateNetworkMappingInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into VMmToAzureUpdateNetworkMappingInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "VmmToVmm") {
		var out VMmToVMmUpdateNetworkMappingInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into VMmToVMmUpdateNetworkMappingInput: %+v", err)
		}
		return out, nil
	}

	var parent BaseFabricSpecificUpdateNetworkMappingInputImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseFabricSpecificUpdateNetworkMappingInputImpl: %+v", err)
	}

	return RawFabricSpecificUpdateNetworkMappingInputImpl{
		fabricSpecificUpdateNetworkMappingInput: parent,
		Type:                                    value,
		Values:                                  temp,
	}, nil

}

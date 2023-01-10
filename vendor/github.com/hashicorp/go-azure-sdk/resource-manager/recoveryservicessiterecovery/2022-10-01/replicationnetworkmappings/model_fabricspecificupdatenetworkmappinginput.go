package replicationnetworkmappings

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FabricSpecificUpdateNetworkMappingInput interface {
}

func unmarshalFabricSpecificUpdateNetworkMappingInputImplementation(input []byte) (FabricSpecificUpdateNetworkMappingInput, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling FabricSpecificUpdateNetworkMappingInput into map[string]interface: %+v", err)
	}

	value, ok := temp["instanceType"].(string)
	if !ok {
		return nil, nil
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

	type RawFabricSpecificUpdateNetworkMappingInputImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawFabricSpecificUpdateNetworkMappingInputImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}

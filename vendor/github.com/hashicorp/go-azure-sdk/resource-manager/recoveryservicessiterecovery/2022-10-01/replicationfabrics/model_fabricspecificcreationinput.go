package replicationfabrics

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FabricSpecificCreationInput interface {
}

// RawModeOfTransitImpl is returned when the Discriminated Value
// doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawFabricSpecificCreationInputImpl struct {
	Type   string
	Values map[string]interface{}
}

func unmarshalFabricSpecificCreationInputImplementation(input []byte) (FabricSpecificCreationInput, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling FabricSpecificCreationInput into map[string]interface: %+v", err)
	}

	value, ok := temp["instanceType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "Azure") {
		var out AzureFabricCreationInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureFabricCreationInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "InMageRcm") {
		var out InMageRcmFabricCreationInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into InMageRcmFabricCreationInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "VMwareV2") {
		var out VMwareV2FabricCreationInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into VMwareV2FabricCreationInput: %+v", err)
		}
		return out, nil
	}

	out := RawFabricSpecificCreationInputImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}

package replicationfabrics

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FabricSpecificCreationInput interface {
	FabricSpecificCreationInput() BaseFabricSpecificCreationInputImpl
}

var _ FabricSpecificCreationInput = BaseFabricSpecificCreationInputImpl{}

type BaseFabricSpecificCreationInputImpl struct {
	InstanceType string `json:"instanceType"`
}

func (s BaseFabricSpecificCreationInputImpl) FabricSpecificCreationInput() BaseFabricSpecificCreationInputImpl {
	return s
}

var _ FabricSpecificCreationInput = RawFabricSpecificCreationInputImpl{}

// RawFabricSpecificCreationInputImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawFabricSpecificCreationInputImpl struct {
	fabricSpecificCreationInput BaseFabricSpecificCreationInputImpl
	Type                        string
	Values                      map[string]interface{}
}

func (s RawFabricSpecificCreationInputImpl) FabricSpecificCreationInput() BaseFabricSpecificCreationInputImpl {
	return s.fabricSpecificCreationInput
}

func UnmarshalFabricSpecificCreationInputImplementation(input []byte) (FabricSpecificCreationInput, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling FabricSpecificCreationInput into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["instanceType"]; ok {
		value = fmt.Sprintf("%v", v)
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

	var parent BaseFabricSpecificCreationInputImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseFabricSpecificCreationInputImpl: %+v", err)
	}

	return RawFabricSpecificCreationInputImpl{
		fabricSpecificCreationInput: parent,
		Type:                        value,
		Values:                      temp,
	}, nil

}

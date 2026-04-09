package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AddDisksProviderSpecificInput interface {
	AddDisksProviderSpecificInput() BaseAddDisksProviderSpecificInputImpl
}

var _ AddDisksProviderSpecificInput = BaseAddDisksProviderSpecificInputImpl{}

type BaseAddDisksProviderSpecificInputImpl struct {
	InstanceType string `json:"instanceType"`
}

func (s BaseAddDisksProviderSpecificInputImpl) AddDisksProviderSpecificInput() BaseAddDisksProviderSpecificInputImpl {
	return s
}

var _ AddDisksProviderSpecificInput = RawAddDisksProviderSpecificInputImpl{}

// RawAddDisksProviderSpecificInputImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawAddDisksProviderSpecificInputImpl struct {
	addDisksProviderSpecificInput BaseAddDisksProviderSpecificInputImpl
	Type                          string
	Values                        map[string]interface{}
}

func (s RawAddDisksProviderSpecificInputImpl) AddDisksProviderSpecificInput() BaseAddDisksProviderSpecificInputImpl {
	return s.addDisksProviderSpecificInput
}

func UnmarshalAddDisksProviderSpecificInputImplementation(input []byte) (AddDisksProviderSpecificInput, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling AddDisksProviderSpecificInput into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["instanceType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "A2A") {
		var out A2AAddDisksInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into A2AAddDisksInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "InMageRcm") {
		var out InMageRcmAddDisksInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into InMageRcmAddDisksInput: %+v", err)
		}
		return out, nil
	}

	var parent BaseAddDisksProviderSpecificInputImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseAddDisksProviderSpecificInputImpl: %+v", err)
	}

	return RawAddDisksProviderSpecificInputImpl{
		addDisksProviderSpecificInput: parent,
		Type:                          value,
		Values:                        temp,
	}, nil

}

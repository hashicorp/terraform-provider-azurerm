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

// RawAddDisksProviderSpecificInputImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawAddDisksProviderSpecificInputImpl struct {
	addDisksProviderSpecificInput BaseAddDisksProviderSpecificInputImpl
	Type                          string
	Values                        map[string]interface{}
}

func (s RawAddDisksProviderSpecificInputImpl) AddDisksProviderSpecificInput() BaseAddDisksProviderSpecificInputImpl {
	return s.addDisksProviderSpecificInput
}

func (s RawAddDisksProviderSpecificInputImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
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

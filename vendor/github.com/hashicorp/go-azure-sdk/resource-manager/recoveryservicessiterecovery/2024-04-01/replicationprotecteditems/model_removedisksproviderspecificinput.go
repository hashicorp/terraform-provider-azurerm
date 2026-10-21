package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RemoveDisksProviderSpecificInput interface {
	RemoveDisksProviderSpecificInput() BaseRemoveDisksProviderSpecificInputImpl
}

var _ RemoveDisksProviderSpecificInput = BaseRemoveDisksProviderSpecificInputImpl{}

type BaseRemoveDisksProviderSpecificInputImpl struct {
	InstanceType string `json:"instanceType"`
}

func (s BaseRemoveDisksProviderSpecificInputImpl) RemoveDisksProviderSpecificInput() BaseRemoveDisksProviderSpecificInputImpl {
	return s
}

var _ RemoveDisksProviderSpecificInput = RawRemoveDisksProviderSpecificInputImpl{}

// RawRemoveDisksProviderSpecificInputImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawRemoveDisksProviderSpecificInputImpl struct {
	removeDisksProviderSpecificInput BaseRemoveDisksProviderSpecificInputImpl
	Type                             string
	Values                           map[string]interface{}
}

func (s RawRemoveDisksProviderSpecificInputImpl) RemoveDisksProviderSpecificInput() BaseRemoveDisksProviderSpecificInputImpl {
	return s.removeDisksProviderSpecificInput
}

func (s RawRemoveDisksProviderSpecificInputImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalRemoveDisksProviderSpecificInputImplementation(input []byte) (RemoveDisksProviderSpecificInput, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling RemoveDisksProviderSpecificInput into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["instanceType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "A2A") {
		var out A2ARemoveDisksInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into A2ARemoveDisksInput: %+v", err)
		}
		return out, nil
	}

	var parent BaseRemoveDisksProviderSpecificInputImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseRemoveDisksProviderSpecificInputImpl: %+v", err)
	}

	return RawRemoveDisksProviderSpecificInputImpl{
		removeDisksProviderSpecificInput: parent,
		Type:                             value,
		Values:                           temp,
	}, nil

}

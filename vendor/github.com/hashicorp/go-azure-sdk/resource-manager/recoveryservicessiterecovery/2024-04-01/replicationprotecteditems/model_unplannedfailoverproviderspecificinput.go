package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UnplannedFailoverProviderSpecificInput interface {
	UnplannedFailoverProviderSpecificInput() BaseUnplannedFailoverProviderSpecificInputImpl
}

var _ UnplannedFailoverProviderSpecificInput = BaseUnplannedFailoverProviderSpecificInputImpl{}

type BaseUnplannedFailoverProviderSpecificInputImpl struct {
	InstanceType string `json:"instanceType"`
}

func (s BaseUnplannedFailoverProviderSpecificInputImpl) UnplannedFailoverProviderSpecificInput() BaseUnplannedFailoverProviderSpecificInputImpl {
	return s
}

var _ UnplannedFailoverProviderSpecificInput = RawUnplannedFailoverProviderSpecificInputImpl{}

// RawUnplannedFailoverProviderSpecificInputImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawUnplannedFailoverProviderSpecificInputImpl struct {
	unplannedFailoverProviderSpecificInput BaseUnplannedFailoverProviderSpecificInputImpl
	Type                                   string
	Values                                 map[string]interface{}
}

func (s RawUnplannedFailoverProviderSpecificInputImpl) UnplannedFailoverProviderSpecificInput() BaseUnplannedFailoverProviderSpecificInputImpl {
	return s.unplannedFailoverProviderSpecificInput
}

func UnmarshalUnplannedFailoverProviderSpecificInputImplementation(input []byte) (UnplannedFailoverProviderSpecificInput, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling UnplannedFailoverProviderSpecificInput into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["instanceType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "A2A") {
		var out A2AUnplannedFailoverInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into A2AUnplannedFailoverInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "HyperVReplicaAzure") {
		var out HyperVReplicaAzureUnplannedFailoverInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into HyperVReplicaAzureUnplannedFailoverInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "InMageAzureV2") {
		var out InMageAzureV2UnplannedFailoverInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into InMageAzureV2UnplannedFailoverInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "InMageRcm") {
		var out InMageRcmUnplannedFailoverInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into InMageRcmUnplannedFailoverInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "InMage") {
		var out InMageUnplannedFailoverInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into InMageUnplannedFailoverInput: %+v", err)
		}
		return out, nil
	}

	var parent BaseUnplannedFailoverProviderSpecificInputImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseUnplannedFailoverProviderSpecificInputImpl: %+v", err)
	}

	return RawUnplannedFailoverProviderSpecificInputImpl{
		unplannedFailoverProviderSpecificInput: parent,
		Type:                                   value,
		Values:                                 temp,
	}, nil

}

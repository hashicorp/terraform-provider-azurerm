package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReverseReplicationProviderSpecificInput interface {
	ReverseReplicationProviderSpecificInput() BaseReverseReplicationProviderSpecificInputImpl
}

var _ ReverseReplicationProviderSpecificInput = BaseReverseReplicationProviderSpecificInputImpl{}

type BaseReverseReplicationProviderSpecificInputImpl struct {
	InstanceType string `json:"instanceType"`
}

func (s BaseReverseReplicationProviderSpecificInputImpl) ReverseReplicationProviderSpecificInput() BaseReverseReplicationProviderSpecificInputImpl {
	return s
}

var _ ReverseReplicationProviderSpecificInput = RawReverseReplicationProviderSpecificInputImpl{}

// RawReverseReplicationProviderSpecificInputImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawReverseReplicationProviderSpecificInputImpl struct {
	reverseReplicationProviderSpecificInput BaseReverseReplicationProviderSpecificInputImpl
	Type                                    string
	Values                                  map[string]interface{}
}

func (s RawReverseReplicationProviderSpecificInputImpl) ReverseReplicationProviderSpecificInput() BaseReverseReplicationProviderSpecificInputImpl {
	return s.reverseReplicationProviderSpecificInput
}

func UnmarshalReverseReplicationProviderSpecificInputImplementation(input []byte) (ReverseReplicationProviderSpecificInput, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ReverseReplicationProviderSpecificInput into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["instanceType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "A2A") {
		var out A2AReprotectInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into A2AReprotectInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "HyperVReplicaAzure") {
		var out HyperVReplicaAzureReprotectInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into HyperVReplicaAzureReprotectInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "InMageAzureV2") {
		var out InMageAzureV2ReprotectInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into InMageAzureV2ReprotectInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "InMageRcmFailback") {
		var out InMageRcmFailbackReprotectInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into InMageRcmFailbackReprotectInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "InMageRcm") {
		var out InMageRcmReprotectInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into InMageRcmReprotectInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "InMage") {
		var out InMageReprotectInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into InMageReprotectInput: %+v", err)
		}
		return out, nil
	}

	var parent BaseReverseReplicationProviderSpecificInputImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseReverseReplicationProviderSpecificInputImpl: %+v", err)
	}

	return RawReverseReplicationProviderSpecificInputImpl{
		reverseReplicationProviderSpecificInput: parent,
		Type:                                    value,
		Values:                                  temp,
	}, nil

}

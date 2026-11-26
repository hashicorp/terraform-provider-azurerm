package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UpdateReplicationProtectedItemProviderInput interface {
	UpdateReplicationProtectedItemProviderInput() BaseUpdateReplicationProtectedItemProviderInputImpl
}

var _ UpdateReplicationProtectedItemProviderInput = BaseUpdateReplicationProtectedItemProviderInputImpl{}

type BaseUpdateReplicationProtectedItemProviderInputImpl struct {
	InstanceType string `json:"instanceType"`
}

func (s BaseUpdateReplicationProtectedItemProviderInputImpl) UpdateReplicationProtectedItemProviderInput() BaseUpdateReplicationProtectedItemProviderInputImpl {
	return s
}

var _ UpdateReplicationProtectedItemProviderInput = RawUpdateReplicationProtectedItemProviderInputImpl{}

// RawUpdateReplicationProtectedItemProviderInputImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawUpdateReplicationProtectedItemProviderInputImpl struct {
	updateReplicationProtectedItemProviderInput BaseUpdateReplicationProtectedItemProviderInputImpl
	Type                                        string
	Values                                      map[string]interface{}
}

func (s RawUpdateReplicationProtectedItemProviderInputImpl) UpdateReplicationProtectedItemProviderInput() BaseUpdateReplicationProtectedItemProviderInputImpl {
	return s.updateReplicationProtectedItemProviderInput
}

func (s RawUpdateReplicationProtectedItemProviderInputImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalUpdateReplicationProtectedItemProviderInputImplementation(input []byte) (UpdateReplicationProtectedItemProviderInput, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling UpdateReplicationProtectedItemProviderInput into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["instanceType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "A2A") {
		var out A2AUpdateReplicationProtectedItemInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into A2AUpdateReplicationProtectedItemInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "HyperVReplicaAzure") {
		var out HyperVReplicaAzureUpdateReplicationProtectedItemInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into HyperVReplicaAzureUpdateReplicationProtectedItemInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "InMageAzureV2") {
		var out InMageAzureV2UpdateReplicationProtectedItemInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into InMageAzureV2UpdateReplicationProtectedItemInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "InMageRcm") {
		var out InMageRcmUpdateReplicationProtectedItemInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into InMageRcmUpdateReplicationProtectedItemInput: %+v", err)
		}
		return out, nil
	}

	var parent BaseUpdateReplicationProtectedItemProviderInputImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseUpdateReplicationProtectedItemProviderInputImpl: %+v", err)
	}

	return RawUpdateReplicationProtectedItemProviderInputImpl{
		updateReplicationProtectedItemProviderInput: parent,
		Type:   value,
		Values: temp,
	}, nil

}

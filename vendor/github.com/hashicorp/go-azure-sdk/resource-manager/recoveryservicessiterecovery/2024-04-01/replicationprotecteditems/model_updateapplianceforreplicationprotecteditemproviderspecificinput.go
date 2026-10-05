package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UpdateApplianceForReplicationProtectedItemProviderSpecificInput interface {
	UpdateApplianceForReplicationProtectedItemProviderSpecificInput() BaseUpdateApplianceForReplicationProtectedItemProviderSpecificInputImpl
}

var _ UpdateApplianceForReplicationProtectedItemProviderSpecificInput = BaseUpdateApplianceForReplicationProtectedItemProviderSpecificInputImpl{}

type BaseUpdateApplianceForReplicationProtectedItemProviderSpecificInputImpl struct {
	InstanceType string `json:"instanceType"`
}

func (s BaseUpdateApplianceForReplicationProtectedItemProviderSpecificInputImpl) UpdateApplianceForReplicationProtectedItemProviderSpecificInput() BaseUpdateApplianceForReplicationProtectedItemProviderSpecificInputImpl {
	return s
}

var _ UpdateApplianceForReplicationProtectedItemProviderSpecificInput = RawUpdateApplianceForReplicationProtectedItemProviderSpecificInputImpl{}

// RawUpdateApplianceForReplicationProtectedItemProviderSpecificInputImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawUpdateApplianceForReplicationProtectedItemProviderSpecificInputImpl struct {
	updateApplianceForReplicationProtectedItemProviderSpecificInput BaseUpdateApplianceForReplicationProtectedItemProviderSpecificInputImpl
	Type                                                            string
	Values                                                          map[string]interface{}
}

func (s RawUpdateApplianceForReplicationProtectedItemProviderSpecificInputImpl) UpdateApplianceForReplicationProtectedItemProviderSpecificInput() BaseUpdateApplianceForReplicationProtectedItemProviderSpecificInputImpl {
	return s.updateApplianceForReplicationProtectedItemProviderSpecificInput
}

func (s RawUpdateApplianceForReplicationProtectedItemProviderSpecificInputImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalUpdateApplianceForReplicationProtectedItemProviderSpecificInputImplementation(input []byte) (UpdateApplianceForReplicationProtectedItemProviderSpecificInput, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling UpdateApplianceForReplicationProtectedItemProviderSpecificInput into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["instanceType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "InMageRcm") {
		var out InMageRcmUpdateApplianceForReplicationProtectedItemInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into InMageRcmUpdateApplianceForReplicationProtectedItemInput: %+v", err)
		}
		return out, nil
	}

	var parent BaseUpdateApplianceForReplicationProtectedItemProviderSpecificInputImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseUpdateApplianceForReplicationProtectedItemProviderSpecificInputImpl: %+v", err)
	}

	return RawUpdateApplianceForReplicationProtectedItemProviderSpecificInputImpl{
		updateApplianceForReplicationProtectedItemProviderSpecificInput: parent,
		Type:   value,
		Values: temp,
	}, nil

}

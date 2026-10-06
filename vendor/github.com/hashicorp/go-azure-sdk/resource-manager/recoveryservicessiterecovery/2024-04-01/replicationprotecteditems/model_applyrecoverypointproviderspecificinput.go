package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplyRecoveryPointProviderSpecificInput interface {
	ApplyRecoveryPointProviderSpecificInput() BaseApplyRecoveryPointProviderSpecificInputImpl
}

var _ ApplyRecoveryPointProviderSpecificInput = BaseApplyRecoveryPointProviderSpecificInputImpl{}

type BaseApplyRecoveryPointProviderSpecificInputImpl struct {
	InstanceType string `json:"instanceType"`
}

func (s BaseApplyRecoveryPointProviderSpecificInputImpl) ApplyRecoveryPointProviderSpecificInput() BaseApplyRecoveryPointProviderSpecificInputImpl {
	return s
}

var _ ApplyRecoveryPointProviderSpecificInput = RawApplyRecoveryPointProviderSpecificInputImpl{}

// RawApplyRecoveryPointProviderSpecificInputImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawApplyRecoveryPointProviderSpecificInputImpl struct {
	applyRecoveryPointProviderSpecificInput BaseApplyRecoveryPointProviderSpecificInputImpl
	Type                                    string
	Values                                  map[string]interface{}
}

func (s RawApplyRecoveryPointProviderSpecificInputImpl) ApplyRecoveryPointProviderSpecificInput() BaseApplyRecoveryPointProviderSpecificInputImpl {
	return s.applyRecoveryPointProviderSpecificInput
}

func (s RawApplyRecoveryPointProviderSpecificInputImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalApplyRecoveryPointProviderSpecificInputImplementation(input []byte) (ApplyRecoveryPointProviderSpecificInput, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ApplyRecoveryPointProviderSpecificInput into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["instanceType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "A2A") {
		var out A2AApplyRecoveryPointInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into A2AApplyRecoveryPointInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "A2ACrossClusterMigration") {
		var out A2ACrossClusterMigrationApplyRecoveryPointInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into A2ACrossClusterMigrationApplyRecoveryPointInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "HyperVReplicaAzure") {
		var out HyperVReplicaAzureApplyRecoveryPointInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into HyperVReplicaAzureApplyRecoveryPointInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "InMageAzureV2") {
		var out InMageAzureV2ApplyRecoveryPointInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into InMageAzureV2ApplyRecoveryPointInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "InMageRcm") {
		var out InMageRcmApplyRecoveryPointInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into InMageRcmApplyRecoveryPointInput: %+v", err)
		}
		return out, nil
	}

	var parent BaseApplyRecoveryPointProviderSpecificInputImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseApplyRecoveryPointProviderSpecificInputImpl: %+v", err)
	}

	return RawApplyRecoveryPointProviderSpecificInputImpl{
		applyRecoveryPointProviderSpecificInput: parent,
		Type:                                    value,
		Values:                                  temp,
	}, nil

}

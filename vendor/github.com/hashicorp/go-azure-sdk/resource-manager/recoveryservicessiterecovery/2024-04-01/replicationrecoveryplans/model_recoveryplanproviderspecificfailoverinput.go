package replicationrecoveryplans

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RecoveryPlanProviderSpecificFailoverInput interface {
	RecoveryPlanProviderSpecificFailoverInput() BaseRecoveryPlanProviderSpecificFailoverInputImpl
}

var _ RecoveryPlanProviderSpecificFailoverInput = BaseRecoveryPlanProviderSpecificFailoverInputImpl{}

type BaseRecoveryPlanProviderSpecificFailoverInputImpl struct {
	InstanceType string `json:"instanceType"`
}

func (s BaseRecoveryPlanProviderSpecificFailoverInputImpl) RecoveryPlanProviderSpecificFailoverInput() BaseRecoveryPlanProviderSpecificFailoverInputImpl {
	return s
}

var _ RecoveryPlanProviderSpecificFailoverInput = RawRecoveryPlanProviderSpecificFailoverInputImpl{}

// RawRecoveryPlanProviderSpecificFailoverInputImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawRecoveryPlanProviderSpecificFailoverInputImpl struct {
	recoveryPlanProviderSpecificFailoverInput BaseRecoveryPlanProviderSpecificFailoverInputImpl
	Type                                      string
	Values                                    map[string]interface{}
}

func (s RawRecoveryPlanProviderSpecificFailoverInputImpl) RecoveryPlanProviderSpecificFailoverInput() BaseRecoveryPlanProviderSpecificFailoverInputImpl {
	return s.recoveryPlanProviderSpecificFailoverInput
}

func UnmarshalRecoveryPlanProviderSpecificFailoverInputImplementation(input []byte) (RecoveryPlanProviderSpecificFailoverInput, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling RecoveryPlanProviderSpecificFailoverInput into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["instanceType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "A2A") {
		var out RecoveryPlanA2AFailoverInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into RecoveryPlanA2AFailoverInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "HyperVReplicaAzureFailback") {
		var out RecoveryPlanHyperVReplicaAzureFailbackInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into RecoveryPlanHyperVReplicaAzureFailbackInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "HyperVReplicaAzure") {
		var out RecoveryPlanHyperVReplicaAzureFailoverInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into RecoveryPlanHyperVReplicaAzureFailoverInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "InMageAzureV2") {
		var out RecoveryPlanInMageAzureV2FailoverInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into RecoveryPlanInMageAzureV2FailoverInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "InMage") {
		var out RecoveryPlanInMageFailoverInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into RecoveryPlanInMageFailoverInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "InMageRcmFailback") {
		var out RecoveryPlanInMageRcmFailbackFailoverInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into RecoveryPlanInMageRcmFailbackFailoverInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "InMageRcm") {
		var out RecoveryPlanInMageRcmFailoverInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into RecoveryPlanInMageRcmFailoverInput: %+v", err)
		}
		return out, nil
	}

	var parent BaseRecoveryPlanProviderSpecificFailoverInputImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseRecoveryPlanProviderSpecificFailoverInputImpl: %+v", err)
	}

	return RawRecoveryPlanProviderSpecificFailoverInputImpl{
		recoveryPlanProviderSpecificFailoverInput: parent,
		Type:   value,
		Values: temp,
	}, nil

}

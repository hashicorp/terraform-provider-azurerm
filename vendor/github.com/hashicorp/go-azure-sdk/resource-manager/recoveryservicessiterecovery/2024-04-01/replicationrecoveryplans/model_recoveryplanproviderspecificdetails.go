package replicationrecoveryplans

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RecoveryPlanProviderSpecificDetails interface {
	RecoveryPlanProviderSpecificDetails() BaseRecoveryPlanProviderSpecificDetailsImpl
}

var _ RecoveryPlanProviderSpecificDetails = BaseRecoveryPlanProviderSpecificDetailsImpl{}

type BaseRecoveryPlanProviderSpecificDetailsImpl struct {
	InstanceType string `json:"instanceType"`
}

func (s BaseRecoveryPlanProviderSpecificDetailsImpl) RecoveryPlanProviderSpecificDetails() BaseRecoveryPlanProviderSpecificDetailsImpl {
	return s
}

var _ RecoveryPlanProviderSpecificDetails = RawRecoveryPlanProviderSpecificDetailsImpl{}

// RawRecoveryPlanProviderSpecificDetailsImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawRecoveryPlanProviderSpecificDetailsImpl struct {
	recoveryPlanProviderSpecificDetails BaseRecoveryPlanProviderSpecificDetailsImpl
	Type                                string
	Values                              map[string]interface{}
}

func (s RawRecoveryPlanProviderSpecificDetailsImpl) RecoveryPlanProviderSpecificDetails() BaseRecoveryPlanProviderSpecificDetailsImpl {
	return s.recoveryPlanProviderSpecificDetails
}

func (s RawRecoveryPlanProviderSpecificDetailsImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalRecoveryPlanProviderSpecificDetailsImplementation(input []byte) (RecoveryPlanProviderSpecificDetails, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling RecoveryPlanProviderSpecificDetails into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["instanceType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "A2A") {
		var out RecoveryPlanA2ADetails
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into RecoveryPlanA2ADetails: %+v", err)
		}
		return out, nil
	}

	var parent BaseRecoveryPlanProviderSpecificDetailsImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseRecoveryPlanProviderSpecificDetailsImpl: %+v", err)
	}

	return RawRecoveryPlanProviderSpecificDetailsImpl{
		recoveryPlanProviderSpecificDetails: parent,
		Type:                                value,
		Values:                              temp,
	}, nil

}

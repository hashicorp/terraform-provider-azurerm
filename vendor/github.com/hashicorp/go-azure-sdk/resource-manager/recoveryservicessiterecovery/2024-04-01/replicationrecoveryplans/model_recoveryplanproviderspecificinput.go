package replicationrecoveryplans

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RecoveryPlanProviderSpecificInput interface {
	RecoveryPlanProviderSpecificInput() BaseRecoveryPlanProviderSpecificInputImpl
}

var _ RecoveryPlanProviderSpecificInput = BaseRecoveryPlanProviderSpecificInputImpl{}

type BaseRecoveryPlanProviderSpecificInputImpl struct {
	InstanceType string `json:"instanceType"`
}

func (s BaseRecoveryPlanProviderSpecificInputImpl) RecoveryPlanProviderSpecificInput() BaseRecoveryPlanProviderSpecificInputImpl {
	return s
}

var _ RecoveryPlanProviderSpecificInput = RawRecoveryPlanProviderSpecificInputImpl{}

// RawRecoveryPlanProviderSpecificInputImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawRecoveryPlanProviderSpecificInputImpl struct {
	recoveryPlanProviderSpecificInput BaseRecoveryPlanProviderSpecificInputImpl
	Type                              string
	Values                            map[string]interface{}
}

func (s RawRecoveryPlanProviderSpecificInputImpl) RecoveryPlanProviderSpecificInput() BaseRecoveryPlanProviderSpecificInputImpl {
	return s.recoveryPlanProviderSpecificInput
}

func UnmarshalRecoveryPlanProviderSpecificInputImplementation(input []byte) (RecoveryPlanProviderSpecificInput, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling RecoveryPlanProviderSpecificInput into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["instanceType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "A2A") {
		var out RecoveryPlanA2AInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into RecoveryPlanA2AInput: %+v", err)
		}
		return out, nil
	}

	var parent BaseRecoveryPlanProviderSpecificInputImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseRecoveryPlanProviderSpecificInputImpl: %+v", err)
	}

	return RawRecoveryPlanProviderSpecificInputImpl{
		recoveryPlanProviderSpecificInput: parent,
		Type:                              value,
		Values:                            temp,
	}, nil

}

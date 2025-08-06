package replicationrecoveryplans

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RecoveryPlanActionDetails interface {
	RecoveryPlanActionDetails() BaseRecoveryPlanActionDetailsImpl
}

var _ RecoveryPlanActionDetails = BaseRecoveryPlanActionDetailsImpl{}

type BaseRecoveryPlanActionDetailsImpl struct {
	InstanceType string `json:"instanceType"`
}

func (s BaseRecoveryPlanActionDetailsImpl) RecoveryPlanActionDetails() BaseRecoveryPlanActionDetailsImpl {
	return s
}

var _ RecoveryPlanActionDetails = RawRecoveryPlanActionDetailsImpl{}

// RawRecoveryPlanActionDetailsImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawRecoveryPlanActionDetailsImpl struct {
	recoveryPlanActionDetails BaseRecoveryPlanActionDetailsImpl
	Type                      string
	Values                    map[string]interface{}
}

func (s RawRecoveryPlanActionDetailsImpl) RecoveryPlanActionDetails() BaseRecoveryPlanActionDetailsImpl {
	return s.recoveryPlanActionDetails
}

func UnmarshalRecoveryPlanActionDetailsImplementation(input []byte) (RecoveryPlanActionDetails, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling RecoveryPlanActionDetails into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["instanceType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "AutomationRunbookActionDetails") {
		var out RecoveryPlanAutomationRunbookActionDetails
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into RecoveryPlanAutomationRunbookActionDetails: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ManualActionDetails") {
		var out RecoveryPlanManualActionDetails
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into RecoveryPlanManualActionDetails: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ScriptActionDetails") {
		var out RecoveryPlanScriptActionDetails
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into RecoveryPlanScriptActionDetails: %+v", err)
		}
		return out, nil
	}

	var parent BaseRecoveryPlanActionDetailsImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseRecoveryPlanActionDetailsImpl: %+v", err)
	}

	return RawRecoveryPlanActionDetailsImpl{
		recoveryPlanActionDetails: parent,
		Type:                      value,
		Values:                    temp,
	}, nil

}

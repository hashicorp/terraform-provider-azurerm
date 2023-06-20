package replicationrecoveryplans

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RecoveryPlanActionDetails interface {
}

func unmarshalRecoveryPlanActionDetailsImplementation(input []byte) (RecoveryPlanActionDetails, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling RecoveryPlanActionDetails into map[string]interface: %+v", err)
	}

	value, ok := temp["instanceType"].(string)
	if !ok {
		return nil, nil
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

	type RawRecoveryPlanActionDetailsImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawRecoveryPlanActionDetailsImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}

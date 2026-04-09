package replicationrecoveryplans

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ RecoveryPlanActionDetails = RecoveryPlanScriptActionDetails{}

type RecoveryPlanScriptActionDetails struct {
	FabricLocation RecoveryPlanActionLocation `json:"fabricLocation"`
	Path           string                     `json:"path"`
	Timeout        *string                    `json:"timeout,omitempty"`

	// Fields inherited from RecoveryPlanActionDetails

	InstanceType string `json:"instanceType"`
}

func (s RecoveryPlanScriptActionDetails) RecoveryPlanActionDetails() BaseRecoveryPlanActionDetailsImpl {
	return BaseRecoveryPlanActionDetailsImpl{
		InstanceType: s.InstanceType,
	}
}

var _ json.Marshaler = RecoveryPlanScriptActionDetails{}

func (s RecoveryPlanScriptActionDetails) MarshalJSON() ([]byte, error) {
	type wrapper RecoveryPlanScriptActionDetails
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling RecoveryPlanScriptActionDetails: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling RecoveryPlanScriptActionDetails: %+v", err)
	}

	decoded["instanceType"] = "ScriptActionDetails"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling RecoveryPlanScriptActionDetails: %+v", err)
	}

	return encoded, nil
}

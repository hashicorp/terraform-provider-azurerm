package replicationrecoveryplans

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ RecoveryPlanActionDetails = RecoveryPlanAutomationRunbookActionDetails{}

type RecoveryPlanAutomationRunbookActionDetails struct {
	FabricLocation RecoveryPlanActionLocation `json:"fabricLocation"`
	RunbookId      *string                    `json:"runbookId,omitempty"`
	Timeout        *string                    `json:"timeout,omitempty"`

	// Fields inherited from RecoveryPlanActionDetails

	InstanceType string `json:"instanceType"`
}

func (s RecoveryPlanAutomationRunbookActionDetails) RecoveryPlanActionDetails() BaseRecoveryPlanActionDetailsImpl {
	return BaseRecoveryPlanActionDetailsImpl{
		InstanceType: s.InstanceType,
	}
}

var _ json.Marshaler = RecoveryPlanAutomationRunbookActionDetails{}

func (s RecoveryPlanAutomationRunbookActionDetails) MarshalJSON() ([]byte, error) {
	type wrapper RecoveryPlanAutomationRunbookActionDetails
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling RecoveryPlanAutomationRunbookActionDetails: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling RecoveryPlanAutomationRunbookActionDetails: %+v", err)
	}

	decoded["instanceType"] = "AutomationRunbookActionDetails"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling RecoveryPlanAutomationRunbookActionDetails: %+v", err)
	}

	return encoded, nil
}

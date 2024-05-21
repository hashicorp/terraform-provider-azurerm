package replicationrecoveryplans

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ RecoveryPlanActionDetails = RecoveryPlanManualActionDetails{}

type RecoveryPlanManualActionDetails struct {
	Description *string `json:"description,omitempty"`

	// Fields inherited from RecoveryPlanActionDetails
}

var _ json.Marshaler = RecoveryPlanManualActionDetails{}

func (s RecoveryPlanManualActionDetails) MarshalJSON() ([]byte, error) {
	type wrapper RecoveryPlanManualActionDetails
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling RecoveryPlanManualActionDetails: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling RecoveryPlanManualActionDetails: %+v", err)
	}
	decoded["instanceType"] = "ManualActionDetails"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling RecoveryPlanManualActionDetails: %+v", err)
	}

	return encoded, nil
}

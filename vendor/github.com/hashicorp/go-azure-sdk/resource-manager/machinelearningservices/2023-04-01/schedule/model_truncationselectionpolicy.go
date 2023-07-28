package schedule

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ EarlyTerminationPolicy = TruncationSelectionPolicy{}

type TruncationSelectionPolicy struct {
	TruncationPercentage *int64 `json:"truncationPercentage,omitempty"`

	// Fields inherited from EarlyTerminationPolicy
	DelayEvaluation    *int64 `json:"delayEvaluation,omitempty"`
	EvaluationInterval *int64 `json:"evaluationInterval,omitempty"`
}

var _ json.Marshaler = TruncationSelectionPolicy{}

func (s TruncationSelectionPolicy) MarshalJSON() ([]byte, error) {
	type wrapper TruncationSelectionPolicy
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling TruncationSelectionPolicy: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling TruncationSelectionPolicy: %+v", err)
	}
	decoded["policyType"] = "TruncationSelection"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling TruncationSelectionPolicy: %+v", err)
	}

	return encoded, nil
}

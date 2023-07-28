package schedule

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ EarlyTerminationPolicy = MedianStoppingPolicy{}

type MedianStoppingPolicy struct {

	// Fields inherited from EarlyTerminationPolicy
	DelayEvaluation    *int64 `json:"delayEvaluation,omitempty"`
	EvaluationInterval *int64 `json:"evaluationInterval,omitempty"`
}

var _ json.Marshaler = MedianStoppingPolicy{}

func (s MedianStoppingPolicy) MarshalJSON() ([]byte, error) {
	type wrapper MedianStoppingPolicy
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MedianStoppingPolicy: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MedianStoppingPolicy: %+v", err)
	}
	decoded["policyType"] = "MedianStopping"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MedianStoppingPolicy: %+v", err)
	}

	return encoded, nil
}

package job

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ EarlyTerminationPolicy = BanditPolicy{}

type BanditPolicy struct {
	SlackAmount *float64 `json:"slackAmount,omitempty"`
	SlackFactor *float64 `json:"slackFactor,omitempty"`

	// Fields inherited from EarlyTerminationPolicy
	DelayEvaluation    *int64 `json:"delayEvaluation,omitempty"`
	EvaluationInterval *int64 `json:"evaluationInterval,omitempty"`
}

var _ json.Marshaler = BanditPolicy{}

func (s BanditPolicy) MarshalJSON() ([]byte, error) {
	type wrapper BanditPolicy
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling BanditPolicy: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling BanditPolicy: %+v", err)
	}
	decoded["policyType"] = "Bandit"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling BanditPolicy: %+v", err)
	}

	return encoded, nil
}

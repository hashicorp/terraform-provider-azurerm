package job

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ JobLimits = SweepJobLimits{}

type SweepJobLimits struct {
	MaxConcurrentTrials *int64  `json:"maxConcurrentTrials,omitempty"`
	MaxTotalTrials      *int64  `json:"maxTotalTrials,omitempty"`
	TrialTimeout        *string `json:"trialTimeout,omitempty"`

	// Fields inherited from JobLimits
	Timeout *string `json:"timeout,omitempty"`
}

var _ json.Marshaler = SweepJobLimits{}

func (s SweepJobLimits) MarshalJSON() ([]byte, error) {
	type wrapper SweepJobLimits
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SweepJobLimits: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SweepJobLimits: %+v", err)
	}
	decoded["jobLimitsType"] = "Sweep"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SweepJobLimits: %+v", err)
	}

	return encoded, nil
}

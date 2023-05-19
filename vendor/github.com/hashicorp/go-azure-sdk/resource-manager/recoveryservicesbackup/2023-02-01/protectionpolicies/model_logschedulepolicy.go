package protectionpolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ SchedulePolicy = LogSchedulePolicy{}

type LogSchedulePolicy struct {
	ScheduleFrequencyInMins *int64 `json:"scheduleFrequencyInMins,omitempty"`

	// Fields inherited from SchedulePolicy
}

var _ json.Marshaler = LogSchedulePolicy{}

func (s LogSchedulePolicy) MarshalJSON() ([]byte, error) {
	type wrapper LogSchedulePolicy
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling LogSchedulePolicy: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling LogSchedulePolicy: %+v", err)
	}
	decoded["schedulePolicyType"] = "LogSchedulePolicy"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling LogSchedulePolicy: %+v", err)
	}

	return encoded, nil
}

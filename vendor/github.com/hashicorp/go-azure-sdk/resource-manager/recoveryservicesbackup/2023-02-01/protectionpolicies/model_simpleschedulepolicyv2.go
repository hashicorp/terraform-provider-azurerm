package protectionpolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ SchedulePolicy = SimpleSchedulePolicyV2{}

type SimpleSchedulePolicyV2 struct {
	DailySchedule        *DailySchedule   `json:"dailySchedule,omitempty"`
	HourlySchedule       *HourlySchedule  `json:"hourlySchedule,omitempty"`
	ScheduleRunFrequency *ScheduleRunType `json:"scheduleRunFrequency,omitempty"`
	WeeklySchedule       *WeeklySchedule  `json:"weeklySchedule,omitempty"`

	// Fields inherited from SchedulePolicy
}

var _ json.Marshaler = SimpleSchedulePolicyV2{}

func (s SimpleSchedulePolicyV2) MarshalJSON() ([]byte, error) {
	type wrapper SimpleSchedulePolicyV2
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SimpleSchedulePolicyV2: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SimpleSchedulePolicyV2: %+v", err)
	}
	decoded["schedulePolicyType"] = "SimpleSchedulePolicyV2"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SimpleSchedulePolicyV2: %+v", err)
	}

	return encoded, nil
}

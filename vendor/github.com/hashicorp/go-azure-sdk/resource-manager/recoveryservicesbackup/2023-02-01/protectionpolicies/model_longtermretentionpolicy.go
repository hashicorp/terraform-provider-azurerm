package protectionpolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ RetentionPolicy = LongTermRetentionPolicy{}

type LongTermRetentionPolicy struct {
	DailySchedule   *DailyRetentionSchedule   `json:"dailySchedule,omitempty"`
	MonthlySchedule *MonthlyRetentionSchedule `json:"monthlySchedule,omitempty"`
	WeeklySchedule  *WeeklyRetentionSchedule  `json:"weeklySchedule,omitempty"`
	YearlySchedule  *YearlyRetentionSchedule  `json:"yearlySchedule,omitempty"`

	// Fields inherited from RetentionPolicy
}

var _ json.Marshaler = LongTermRetentionPolicy{}

func (s LongTermRetentionPolicy) MarshalJSON() ([]byte, error) {
	type wrapper LongTermRetentionPolicy
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling LongTermRetentionPolicy: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling LongTermRetentionPolicy: %+v", err)
	}
	decoded["retentionPolicyType"] = "LongTermRetentionPolicy"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling LongTermRetentionPolicy: %+v", err)
	}

	return encoded, nil
}

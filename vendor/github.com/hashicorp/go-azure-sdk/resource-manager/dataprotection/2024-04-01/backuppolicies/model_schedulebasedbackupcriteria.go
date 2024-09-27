package backuppolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ BackupCriteria = ScheduleBasedBackupCriteria{}

type ScheduleBasedBackupCriteria struct {
	AbsoluteCriteria *[]AbsoluteMarker `json:"absoluteCriteria,omitempty"`
	DaysOfMonth      *[]Day            `json:"daysOfMonth,omitempty"`
	DaysOfTheWeek    *[]DayOfWeek      `json:"daysOfTheWeek,omitempty"`
	MonthsOfYear     *[]Month          `json:"monthsOfYear,omitempty"`
	ScheduleTimes    *[]string         `json:"scheduleTimes,omitempty"`
	WeeksOfTheMonth  *[]WeekNumber     `json:"weeksOfTheMonth,omitempty"`

	// Fields inherited from BackupCriteria

	ObjectType string `json:"objectType"`
}

func (s ScheduleBasedBackupCriteria) BackupCriteria() BaseBackupCriteriaImpl {
	return BaseBackupCriteriaImpl{
		ObjectType: s.ObjectType,
	}
}

var _ json.Marshaler = ScheduleBasedBackupCriteria{}

func (s ScheduleBasedBackupCriteria) MarshalJSON() ([]byte, error) {
	type wrapper ScheduleBasedBackupCriteria
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ScheduleBasedBackupCriteria: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ScheduleBasedBackupCriteria: %+v", err)
	}

	decoded["objectType"] = "ScheduleBasedBackupCriteria"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ScheduleBasedBackupCriteria: %+v", err)
	}

	return encoded, nil
}

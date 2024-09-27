package appplatform

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ MaintenanceScheduleConfiguration = WeeklyMaintenanceScheduleConfiguration{}

type WeeklyMaintenanceScheduleConfiguration struct {
	Day      WeekDay `json:"day"`
	Duration *string `json:"duration,omitempty"`
	Hour     int64   `json:"hour"`

	// Fields inherited from MaintenanceScheduleConfiguration

	Frequency Frequency `json:"frequency"`
}

func (s WeeklyMaintenanceScheduleConfiguration) MaintenanceScheduleConfiguration() BaseMaintenanceScheduleConfigurationImpl {
	return BaseMaintenanceScheduleConfigurationImpl{
		Frequency: s.Frequency,
	}
}

var _ json.Marshaler = WeeklyMaintenanceScheduleConfiguration{}

func (s WeeklyMaintenanceScheduleConfiguration) MarshalJSON() ([]byte, error) {
	type wrapper WeeklyMaintenanceScheduleConfiguration
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling WeeklyMaintenanceScheduleConfiguration: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling WeeklyMaintenanceScheduleConfiguration: %+v", err)
	}

	decoded["frequency"] = "Weekly"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling WeeklyMaintenanceScheduleConfiguration: %+v", err)
	}

	return encoded, nil
}

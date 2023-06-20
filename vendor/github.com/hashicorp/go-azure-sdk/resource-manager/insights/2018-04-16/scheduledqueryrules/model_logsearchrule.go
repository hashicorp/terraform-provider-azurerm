package scheduledqueryrules

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LogSearchRule struct {
	Action                   Action             `json:"action"`
	AutoMitigate             *bool              `json:"autoMitigate,omitempty"`
	CreatedWithApiVersion    *string            `json:"createdWithApiVersion,omitempty"`
	Description              *string            `json:"description,omitempty"`
	DisplayName              *string            `json:"displayName,omitempty"`
	Enabled                  *Enabled           `json:"enabled,omitempty"`
	IsLegacyLogAnalyticsRule *bool              `json:"isLegacyLogAnalyticsRule,omitempty"`
	LastUpdatedTime          *string            `json:"lastUpdatedTime,omitempty"`
	ProvisioningState        *ProvisioningState `json:"provisioningState,omitempty"`
	Schedule                 *Schedule          `json:"schedule,omitempty"`
	Source                   Source             `json:"source"`
}

func (o *LogSearchRule) GetLastUpdatedTimeAsTime() (*time.Time, error) {
	if o.LastUpdatedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastUpdatedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *LogSearchRule) SetLastUpdatedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastUpdatedTime = &formatted
}

var _ json.Unmarshaler = &LogSearchRule{}

func (s *LogSearchRule) UnmarshalJSON(bytes []byte) error {
	type alias LogSearchRule
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into LogSearchRule: %+v", err)
	}

	s.AutoMitigate = decoded.AutoMitigate
	s.CreatedWithApiVersion = decoded.CreatedWithApiVersion
	s.Description = decoded.Description
	s.DisplayName = decoded.DisplayName
	s.Enabled = decoded.Enabled
	s.IsLegacyLogAnalyticsRule = decoded.IsLegacyLogAnalyticsRule
	s.LastUpdatedTime = decoded.LastUpdatedTime
	s.ProvisioningState = decoded.ProvisioningState
	s.Schedule = decoded.Schedule
	s.Source = decoded.Source

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling LogSearchRule into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["action"]; ok {
		impl, err := unmarshalActionImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Action' for 'LogSearchRule': %+v", err)
		}
		s.Action = impl
	}
	return nil
}

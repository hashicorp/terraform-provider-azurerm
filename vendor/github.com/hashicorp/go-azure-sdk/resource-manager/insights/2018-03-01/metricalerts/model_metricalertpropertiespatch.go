package metricalerts

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MetricAlertPropertiesPatch struct {
	Actions              *[]MetricAlertAction `json:"actions,omitempty"`
	AutoMitigate         *bool                `json:"autoMitigate,omitempty"`
	Criteria             MetricAlertCriteria  `json:"criteria"`
	Description          *string              `json:"description,omitempty"`
	Enabled              *bool                `json:"enabled,omitempty"`
	EvaluationFrequency  *string              `json:"evaluationFrequency,omitempty"`
	IsMigrated           *bool                `json:"isMigrated,omitempty"`
	LastUpdatedTime      *string              `json:"lastUpdatedTime,omitempty"`
	Scopes               *[]string            `json:"scopes,omitempty"`
	Severity             *int64               `json:"severity,omitempty"`
	TargetResourceRegion *string              `json:"targetResourceRegion,omitempty"`
	TargetResourceType   *string              `json:"targetResourceType,omitempty"`
	WindowSize           *string              `json:"windowSize,omitempty"`
}

func (o *MetricAlertPropertiesPatch) GetLastUpdatedTimeAsTime() (*time.Time, error) {
	if o.LastUpdatedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastUpdatedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *MetricAlertPropertiesPatch) SetLastUpdatedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastUpdatedTime = &formatted
}

var _ json.Unmarshaler = &MetricAlertPropertiesPatch{}

func (s *MetricAlertPropertiesPatch) UnmarshalJSON(bytes []byte) error {
	type alias MetricAlertPropertiesPatch
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into MetricAlertPropertiesPatch: %+v", err)
	}

	s.Actions = decoded.Actions
	s.AutoMitigate = decoded.AutoMitigate
	s.Description = decoded.Description
	s.Enabled = decoded.Enabled
	s.EvaluationFrequency = decoded.EvaluationFrequency
	s.IsMigrated = decoded.IsMigrated
	s.LastUpdatedTime = decoded.LastUpdatedTime
	s.Scopes = decoded.Scopes
	s.Severity = decoded.Severity
	s.TargetResourceRegion = decoded.TargetResourceRegion
	s.TargetResourceType = decoded.TargetResourceType
	s.WindowSize = decoded.WindowSize

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling MetricAlertPropertiesPatch into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["criteria"]; ok {
		impl, err := unmarshalMetricAlertCriteriaImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Criteria' for 'MetricAlertPropertiesPatch': %+v", err)
		}
		s.Criteria = impl
	}
	return nil
}

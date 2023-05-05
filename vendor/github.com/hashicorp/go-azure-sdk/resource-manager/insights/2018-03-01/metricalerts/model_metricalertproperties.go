package metricalerts

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MetricAlertProperties struct {
	Actions              *[]MetricAlertAction `json:"actions,omitempty"`
	AutoMitigate         *bool                `json:"autoMitigate,omitempty"`
	Criteria             MetricAlertCriteria  `json:"criteria"`
	Description          *string              `json:"description,omitempty"`
	Enabled              bool                 `json:"enabled"`
	EvaluationFrequency  string               `json:"evaluationFrequency"`
	IsMigrated           *bool                `json:"isMigrated,omitempty"`
	LastUpdatedTime      *string              `json:"lastUpdatedTime,omitempty"`
	Scopes               []string             `json:"scopes"`
	Severity             int64                `json:"severity"`
	TargetResourceRegion *string              `json:"targetResourceRegion,omitempty"`
	TargetResourceType   *string              `json:"targetResourceType,omitempty"`
	WindowSize           string               `json:"windowSize"`
}

func (o *MetricAlertProperties) GetLastUpdatedTimeAsTime() (*time.Time, error) {
	if o.LastUpdatedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastUpdatedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *MetricAlertProperties) SetLastUpdatedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastUpdatedTime = &formatted
}

var _ json.Unmarshaler = &MetricAlertProperties{}

func (s *MetricAlertProperties) UnmarshalJSON(bytes []byte) error {
	type alias MetricAlertProperties
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into MetricAlertProperties: %+v", err)
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
		return fmt.Errorf("unmarshaling MetricAlertProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["criteria"]; ok {
		impl, err := unmarshalMetricAlertCriteriaImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Criteria' for 'MetricAlertProperties': %+v", err)
		}
		s.Criteria = impl
	}
	return nil
}

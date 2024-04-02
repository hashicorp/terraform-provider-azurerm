package alertsmanagements

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Essentials struct {
	ActionStatus                     *ActionStatus     `json:"actionStatus,omitempty"`
	AlertRule                        *string           `json:"alertRule,omitempty"`
	AlertState                       *AlertState       `json:"alertState,omitempty"`
	Description                      *string           `json:"description,omitempty"`
	LastModifiedDateTime             *string           `json:"lastModifiedDateTime,omitempty"`
	LastModifiedUserName             *string           `json:"lastModifiedUserName,omitempty"`
	MonitorCondition                 *MonitorCondition `json:"monitorCondition,omitempty"`
	MonitorConditionResolvedDateTime *string           `json:"monitorConditionResolvedDateTime,omitempty"`
	MonitorService                   *MonitorService   `json:"monitorService,omitempty"`
	Severity                         *Severity         `json:"severity,omitempty"`
	SignalType                       *SignalType       `json:"signalType,omitempty"`
	SmartGroupId                     *string           `json:"smartGroupId,omitempty"`
	SmartGroupingReason              *string           `json:"smartGroupingReason,omitempty"`
	SourceCreatedId                  *string           `json:"sourceCreatedId,omitempty"`
	StartDateTime                    *string           `json:"startDateTime,omitempty"`
	TargetResource                   *string           `json:"targetResource,omitempty"`
	TargetResourceGroup              *string           `json:"targetResourceGroup,omitempty"`
	TargetResourceName               *string           `json:"targetResourceName,omitempty"`
	TargetResourceType               *string           `json:"targetResourceType,omitempty"`
}

func (o *Essentials) GetLastModifiedDateTimeAsTime() (*time.Time, error) {
	if o.LastModifiedDateTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastModifiedDateTime, "2006-01-02T15:04:05Z07:00")
}

func (o *Essentials) SetLastModifiedDateTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModifiedDateTime = &formatted
}

func (o *Essentials) GetMonitorConditionResolvedDateTimeAsTime() (*time.Time, error) {
	if o.MonitorConditionResolvedDateTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.MonitorConditionResolvedDateTime, "2006-01-02T15:04:05Z07:00")
}

func (o *Essentials) SetMonitorConditionResolvedDateTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.MonitorConditionResolvedDateTime = &formatted
}

func (o *Essentials) GetStartDateTimeAsTime() (*time.Time, error) {
	if o.StartDateTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartDateTime, "2006-01-02T15:04:05Z07:00")
}

func (o *Essentials) SetStartDateTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartDateTime = &formatted
}

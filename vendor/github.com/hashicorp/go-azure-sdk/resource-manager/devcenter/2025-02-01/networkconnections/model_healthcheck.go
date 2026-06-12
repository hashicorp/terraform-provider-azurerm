package networkconnections

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HealthCheck struct {
	AdditionalDetails *string            `json:"additionalDetails,omitempty"`
	DisplayName       *string            `json:"displayName,omitempty"`
	EndDateTime       *string            `json:"endDateTime,omitempty"`
	ErrorType         *string            `json:"errorType,omitempty"`
	RecommendedAction *string            `json:"recommendedAction,omitempty"`
	StartDateTime     *string            `json:"startDateTime,omitempty"`
	Status            *HealthCheckStatus `json:"status,omitempty"`
}

func (o *HealthCheck) GetEndDateTimeAsTime() (*time.Time, error) {
	if o.EndDateTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndDateTime, "2006-01-02T15:04:05Z07:00")
}

func (o *HealthCheck) SetEndDateTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndDateTime = &formatted
}

func (o *HealthCheck) GetStartDateTimeAsTime() (*time.Time, error) {
	if o.StartDateTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartDateTime, "2006-01-02T15:04:05Z07:00")
}

func (o *HealthCheck) SetStartDateTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartDateTime = &formatted
}

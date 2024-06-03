package serverendpointresource

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerEndpointProvisioningStepStatus struct {
	AdditionalInformation *map[string]string `json:"additionalInformation,omitempty"`
	EndTime               *string            `json:"endTime,omitempty"`
	ErrorCode             *int64             `json:"errorCode,omitempty"`
	MinutesLeft           *int64             `json:"minutesLeft,omitempty"`
	Name                  *string            `json:"name,omitempty"`
	ProgressPercentage    *int64             `json:"progressPercentage,omitempty"`
	StartTime             *string            `json:"startTime,omitempty"`
	Status                *string            `json:"status,omitempty"`
}

func (o *ServerEndpointProvisioningStepStatus) GetEndTimeAsTime() (*time.Time, error) {
	if o.EndTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ServerEndpointProvisioningStepStatus) SetEndTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndTime = &formatted
}

func (o *ServerEndpointProvisioningStepStatus) GetStartTimeAsTime() (*time.Time, error) {
	if o.StartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ServerEndpointProvisioningStepStatus) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = &formatted
}

package webapps

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeploymentProperties struct {
	Active      *bool   `json:"active,omitempty"`
	Author      *string `json:"author,omitempty"`
	AuthorEmail *string `json:"author_email,omitempty"`
	Deployer    *string `json:"deployer,omitempty"`
	Details     *string `json:"details,omitempty"`
	EndTime     *string `json:"end_time,omitempty"`
	Message     *string `json:"message,omitempty"`
	StartTime   *string `json:"start_time,omitempty"`
	Status      *int64  `json:"status,omitempty"`
}

func (o *DeploymentProperties) GetEndTimeAsTime() (*time.Time, error) {
	if o.EndTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndTime, "2006-01-02T15:04:05Z07:00")
}

func (o *DeploymentProperties) SetEndTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndTime = &formatted
}

func (o *DeploymentProperties) GetStartTimeAsTime() (*time.Time, error) {
	if o.StartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *DeploymentProperties) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = &formatted
}

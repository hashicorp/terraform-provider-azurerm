package containerappsrevisions

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RevisionProperties struct {
	Active            *bool                      `json:"active,omitempty"`
	CreatedTime       *string                    `json:"createdTime,omitempty"`
	Fqdn              *string                    `json:"fqdn,omitempty"`
	HealthState       *RevisionHealthState       `json:"healthState,omitempty"`
	LastActiveTime    *string                    `json:"lastActiveTime,omitempty"`
	ProvisioningError *string                    `json:"provisioningError,omitempty"`
	ProvisioningState *RevisionProvisioningState `json:"provisioningState,omitempty"`
	Replicas          *int64                     `json:"replicas,omitempty"`
	RunningState      *RevisionRunningState      `json:"runningState,omitempty"`
	Template          *Template                  `json:"template,omitempty"`
	TrafficWeight     *int64                     `json:"trafficWeight,omitempty"`
}

func (o *RevisionProperties) GetCreatedTimeAsTime() (*time.Time, error) {
	if o.CreatedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *RevisionProperties) SetCreatedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedTime = &formatted
}

func (o *RevisionProperties) GetLastActiveTimeAsTime() (*time.Time, error) {
	if o.LastActiveTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastActiveTime, "2006-01-02T15:04:05Z07:00")
}

func (o *RevisionProperties) SetLastActiveTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastActiveTime = &formatted
}

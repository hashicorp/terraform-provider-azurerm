package sessionhost

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SessionHostProperties struct {
	AgentVersion                  *string                         `json:"agentVersion,omitempty"`
	AllowNewSession               *bool                           `json:"allowNewSession,omitempty"`
	AssignedUser                  *string                         `json:"assignedUser,omitempty"`
	FriendlyName                  *string                         `json:"friendlyName,omitempty"`
	LastHeartBeat                 *string                         `json:"lastHeartBeat,omitempty"`
	LastUpdateTime                *string                         `json:"lastUpdateTime,omitempty"`
	ObjectId                      *string                         `json:"objectId,omitempty"`
	OsVersion                     *string                         `json:"osVersion,omitempty"`
	ResourceId                    *string                         `json:"resourceId,omitempty"`
	SessionHostHealthCheckResults *[]SessionHostHealthCheckReport `json:"sessionHostHealthCheckResults,omitempty"`
	Sessions                      *int64                          `json:"sessions,omitempty"`
	Status                        *Status                         `json:"status,omitempty"`
	StatusTimestamp               *string                         `json:"statusTimestamp,omitempty"`
	SxSStackVersion               *string                         `json:"sxSStackVersion,omitempty"`
	UpdateErrorMessage            *string                         `json:"updateErrorMessage,omitempty"`
	UpdateState                   *UpdateState                    `json:"updateState,omitempty"`
	VirtualMachineId              *string                         `json:"virtualMachineId,omitempty"`
}

func (o *SessionHostProperties) GetLastHeartBeatAsTime() (*time.Time, error) {
	if o.LastHeartBeat == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastHeartBeat, "2006-01-02T15:04:05Z07:00")
}

func (o *SessionHostProperties) SetLastHeartBeatAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastHeartBeat = &formatted
}

func (o *SessionHostProperties) GetLastUpdateTimeAsTime() (*time.Time, error) {
	if o.LastUpdateTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastUpdateTime, "2006-01-02T15:04:05Z07:00")
}

func (o *SessionHostProperties) SetLastUpdateTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastUpdateTime = &formatted
}

func (o *SessionHostProperties) GetStatusTimestampAsTime() (*time.Time, error) {
	if o.StatusTimestamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StatusTimestamp, "2006-01-02T15:04:05Z07:00")
}

func (o *SessionHostProperties) SetStatusTimestampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StatusTimestamp = &formatted
}

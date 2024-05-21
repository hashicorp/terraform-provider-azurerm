package replicationfabrics

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MasterTargetServer struct {
	AgentExpiryDate         *string            `json:"agentExpiryDate,omitempty"`
	AgentVersion            *string            `json:"agentVersion,omitempty"`
	AgentVersionDetails     *VersionDetails    `json:"agentVersionDetails,omitempty"`
	DataStores              *[]DataStore       `json:"dataStores,omitempty"`
	DiskCount               *int64             `json:"diskCount,omitempty"`
	HealthErrors            *[]HealthError     `json:"healthErrors,omitempty"`
	IPAddress               *string            `json:"ipAddress,omitempty"`
	Id                      *string            `json:"id,omitempty"`
	LastHeartbeat           *string            `json:"lastHeartbeat,omitempty"`
	MarsAgentExpiryDate     *string            `json:"marsAgentExpiryDate,omitempty"`
	MarsAgentVersion        *string            `json:"marsAgentVersion,omitempty"`
	MarsAgentVersionDetails *VersionDetails    `json:"marsAgentVersionDetails,omitempty"`
	Name                    *string            `json:"name,omitempty"`
	OsType                  *string            `json:"osType,omitempty"`
	OsVersion               *string            `json:"osVersion,omitempty"`
	RetentionVolumes        *[]RetentionVolume `json:"retentionVolumes,omitempty"`
	ValidationErrors        *[]HealthError     `json:"validationErrors,omitempty"`
	VersionStatus           *string            `json:"versionStatus,omitempty"`
}

func (o *MasterTargetServer) GetAgentExpiryDateAsTime() (*time.Time, error) {
	if o.AgentExpiryDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.AgentExpiryDate, "2006-01-02T15:04:05Z07:00")
}

func (o *MasterTargetServer) SetAgentExpiryDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.AgentExpiryDate = &formatted
}

func (o *MasterTargetServer) GetLastHeartbeatAsTime() (*time.Time, error) {
	if o.LastHeartbeat == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastHeartbeat, "2006-01-02T15:04:05Z07:00")
}

func (o *MasterTargetServer) SetLastHeartbeatAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastHeartbeat = &formatted
}

func (o *MasterTargetServer) GetMarsAgentExpiryDateAsTime() (*time.Time, error) {
	if o.MarsAgentExpiryDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.MarsAgentExpiryDate, "2006-01-02T15:04:05Z07:00")
}

func (o *MasterTargetServer) SetMarsAgentExpiryDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.MarsAgentExpiryDate = &formatted
}

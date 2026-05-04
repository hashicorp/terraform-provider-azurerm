package replicationfabrics

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProcessServer struct {
	AgentExpiryDate                    *string                  `json:"agentExpiryDate,omitempty"`
	AgentVersion                       *string                  `json:"agentVersion,omitempty"`
	AgentVersionDetails                *VersionDetails          `json:"agentVersionDetails,omitempty"`
	AvailableMemoryInBytes             *int64                   `json:"availableMemoryInBytes,omitempty"`
	AvailableSpaceInBytes              *int64                   `json:"availableSpaceInBytes,omitempty"`
	CpuLoad                            *string                  `json:"cpuLoad,omitempty"`
	CpuLoadStatus                      *string                  `json:"cpuLoadStatus,omitempty"`
	FriendlyName                       *string                  `json:"friendlyName,omitempty"`
	Health                             *ProtectionHealth        `json:"health,omitempty"`
	HealthErrors                       *[]HealthError           `json:"healthErrors,omitempty"`
	HostId                             *string                  `json:"hostId,omitempty"`
	IPAddress                          *string                  `json:"ipAddress,omitempty"`
	Id                                 *string                  `json:"id,omitempty"`
	LastHeartbeat                      *string                  `json:"lastHeartbeat,omitempty"`
	MachineCount                       *string                  `json:"machineCount,omitempty"`
	MarsCommunicationStatus            *string                  `json:"marsCommunicationStatus,omitempty"`
	MarsRegistrationStatus             *string                  `json:"marsRegistrationStatus,omitempty"`
	MemoryUsageStatus                  *string                  `json:"memoryUsageStatus,omitempty"`
	MobilityServiceUpdates             *[]MobilityServiceUpdate `json:"mobilityServiceUpdates,omitempty"`
	OsType                             *string                  `json:"osType,omitempty"`
	OsVersion                          *string                  `json:"osVersion,omitempty"`
	PsServiceStatus                    *string                  `json:"psServiceStatus,omitempty"`
	PsStatsRefreshTime                 *string                  `json:"psStatsRefreshTime,omitempty"`
	ReplicationPairCount               *string                  `json:"replicationPairCount,omitempty"`
	SpaceUsageStatus                   *string                  `json:"spaceUsageStatus,omitempty"`
	SslCertExpiryDate                  *string                  `json:"sslCertExpiryDate,omitempty"`
	SslCertExpiryRemainingDays         *int64                   `json:"sslCertExpiryRemainingDays,omitempty"`
	SystemLoad                         *string                  `json:"systemLoad,omitempty"`
	SystemLoadStatus                   *string                  `json:"systemLoadStatus,omitempty"`
	ThroughputInBytes                  *int64                   `json:"throughputInBytes,omitempty"`
	ThroughputInMBps                   *int64                   `json:"throughputInMBps,omitempty"`
	ThroughputStatus                   *string                  `json:"throughputStatus,omitempty"`
	ThroughputUploadPendingDataInBytes *int64                   `json:"throughputUploadPendingDataInBytes,omitempty"`
	TotalMemoryInBytes                 *int64                   `json:"totalMemoryInBytes,omitempty"`
	TotalSpaceInBytes                  *int64                   `json:"totalSpaceInBytes,omitempty"`
	VersionStatus                      *string                  `json:"versionStatus,omitempty"`
}

func (o *ProcessServer) GetAgentExpiryDateAsTime() (*time.Time, error) {
	if o.AgentExpiryDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.AgentExpiryDate, "2006-01-02T15:04:05Z07:00")
}

func (o *ProcessServer) SetAgentExpiryDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.AgentExpiryDate = &formatted
}

func (o *ProcessServer) GetLastHeartbeatAsTime() (*time.Time, error) {
	if o.LastHeartbeat == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastHeartbeat, "2006-01-02T15:04:05Z07:00")
}

func (o *ProcessServer) SetLastHeartbeatAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastHeartbeat = &formatted
}

func (o *ProcessServer) GetPsStatsRefreshTimeAsTime() (*time.Time, error) {
	if o.PsStatsRefreshTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.PsStatsRefreshTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ProcessServer) SetPsStatsRefreshTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.PsStatsRefreshTime = &formatted
}

func (o *ProcessServer) GetSslCertExpiryDateAsTime() (*time.Time, error) {
	if o.SslCertExpiryDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.SslCertExpiryDate, "2006-01-02T15:04:05Z07:00")
}

func (o *ProcessServer) SetSslCertExpiryDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.SslCertExpiryDate = &formatted
}

package replicationfabrics

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProcessServerDetails struct {
	AvailableMemoryInBytes             *int64              `json:"availableMemoryInBytes,omitempty"`
	AvailableSpaceInBytes              *int64              `json:"availableSpaceInBytes,omitempty"`
	BiosId                             *string             `json:"biosId,omitempty"`
	DiskUsageStatus                    *RcmComponentStatus `json:"diskUsageStatus,omitempty"`
	FabricObjectId                     *string             `json:"fabricObjectId,omitempty"`
	Fqdn                               *string             `json:"fqdn,omitempty"`
	FreeSpacePercentage                *float64            `json:"freeSpacePercentage,omitempty"`
	Health                             *ProtectionHealth   `json:"health,omitempty"`
	HealthErrors                       *[]HealthError      `json:"healthErrors,omitempty"`
	HistoricHealth                     *ProtectionHealth   `json:"historicHealth,omitempty"`
	IPAddresses                        *[]string           `json:"ipAddresses,omitempty"`
	Id                                 *string             `json:"id,omitempty"`
	LastHeartbeatUtc                   *string             `json:"lastHeartbeatUtc,omitempty"`
	MemoryUsagePercentage              *float64            `json:"memoryUsagePercentage,omitempty"`
	MemoryUsageStatus                  *RcmComponentStatus `json:"memoryUsageStatus,omitempty"`
	Name                               *string             `json:"name,omitempty"`
	ProcessorUsagePercentage           *float64            `json:"processorUsagePercentage,omitempty"`
	ProcessorUsageStatus               *RcmComponentStatus `json:"processorUsageStatus,omitempty"`
	ProtectedItemCount                 *int64              `json:"protectedItemCount,omitempty"`
	SystemLoad                         *int64              `json:"systemLoad,omitempty"`
	SystemLoadStatus                   *RcmComponentStatus `json:"systemLoadStatus,omitempty"`
	ThroughputInBytes                  *int64              `json:"throughputInBytes,omitempty"`
	ThroughputStatus                   *RcmComponentStatus `json:"throughputStatus,omitempty"`
	ThroughputUploadPendingDataInBytes *int64              `json:"throughputUploadPendingDataInBytes,omitempty"`
	TotalMemoryInBytes                 *int64              `json:"totalMemoryInBytes,omitempty"`
	TotalSpaceInBytes                  *int64              `json:"totalSpaceInBytes,omitempty"`
	UsedMemoryInBytes                  *int64              `json:"usedMemoryInBytes,omitempty"`
	UsedSpaceInBytes                   *int64              `json:"usedSpaceInBytes,omitempty"`
	Version                            *string             `json:"version,omitempty"`
}

func (o *ProcessServerDetails) GetLastHeartbeatUtcAsTime() (*time.Time, error) {
	if o.LastHeartbeatUtc == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastHeartbeatUtc, "2006-01-02T15:04:05Z07:00")
}

func (o *ProcessServerDetails) SetLastHeartbeatUtcAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastHeartbeatUtc = &formatted
}

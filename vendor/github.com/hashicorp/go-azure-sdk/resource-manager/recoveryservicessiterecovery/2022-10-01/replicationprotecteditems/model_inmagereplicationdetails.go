package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ReplicationProviderSpecificSettings = InMageReplicationDetails{}

type InMageReplicationDetails struct {
	ActiveSiteType               *string                       `json:"activeSiteType,omitempty"`
	AgentDetails                 *InMageAgentDetails           `json:"agentDetails,omitempty"`
	AzureStorageAccountId        *string                       `json:"azureStorageAccountId,omitempty"`
	CompressedDataRateInMB       *float64                      `json:"compressedDataRateInMB,omitempty"`
	ConsistencyPoints            *map[string]string            `json:"consistencyPoints,omitempty"`
	DataStores                   *[]string                     `json:"datastores,omitempty"`
	DiscoveryType                *string                       `json:"discoveryType,omitempty"`
	DiskResized                  *string                       `json:"diskResized,omitempty"`
	IPAddress                    *string                       `json:"ipAddress,omitempty"`
	InfrastructureVMId           *string                       `json:"infrastructureVmId,omitempty"`
	IsAdditionalStatsAvailable   *bool                         `json:"isAdditionalStatsAvailable,omitempty"`
	LastHeartbeat                *string                       `json:"lastHeartbeat,omitempty"`
	LastRpoCalculatedTime        *string                       `json:"lastRpoCalculatedTime,omitempty"`
	LastUpdateReceivedTime       *string                       `json:"lastUpdateReceivedTime,omitempty"`
	MasterTargetId               *string                       `json:"masterTargetId,omitempty"`
	MultiVMGroupId               *string                       `json:"multiVmGroupId,omitempty"`
	MultiVMGroupName             *string                       `json:"multiVmGroupName,omitempty"`
	MultiVMSyncStatus            *string                       `json:"multiVmSyncStatus,omitempty"`
	OsDetails                    *OSDiskDetails                `json:"osDetails,omitempty"`
	OsVersion                    *string                       `json:"osVersion,omitempty"`
	ProcessServerId              *string                       `json:"processServerId,omitempty"`
	ProtectedDisks               *[]InMageProtectedDiskDetails `json:"protectedDisks,omitempty"`
	ProtectionStage              *string                       `json:"protectionStage,omitempty"`
	RebootAfterUpdateStatus      *string                       `json:"rebootAfterUpdateStatus,omitempty"`
	ReplicaId                    *string                       `json:"replicaId,omitempty"`
	ResyncDetails                *InitialReplicationDetails    `json:"resyncDetails,omitempty"`
	RetentionWindowEnd           *string                       `json:"retentionWindowEnd,omitempty"`
	RetentionWindowStart         *string                       `json:"retentionWindowStart,omitempty"`
	RpoInSeconds                 *int64                        `json:"rpoInSeconds,omitempty"`
	SourceVMCPUCount             *int64                        `json:"sourceVmCpuCount,omitempty"`
	SourceVMRamSizeInMB          *int64                        `json:"sourceVmRamSizeInMB,omitempty"`
	TotalDataTransferred         *int64                        `json:"totalDataTransferred,omitempty"`
	TotalProgressHealth          *string                       `json:"totalProgressHealth,omitempty"`
	UncompressedDataRateInMB     *float64                      `json:"uncompressedDataRateInMB,omitempty"`
	VCenterInfrastructureId      *string                       `json:"vCenterInfrastructureId,omitempty"`
	VMId                         *string                       `json:"vmId,omitempty"`
	VMNics                       *[]VMNicDetails               `json:"vmNics,omitempty"`
	VMProtectionState            *string                       `json:"vmProtectionState,omitempty"`
	VMProtectionStateDescription *string                       `json:"vmProtectionStateDescription,omitempty"`
	ValidationErrors             *[]HealthError                `json:"validationErrors,omitempty"`

	// Fields inherited from ReplicationProviderSpecificSettings
}

var _ json.Marshaler = InMageReplicationDetails{}

func (s InMageReplicationDetails) MarshalJSON() ([]byte, error) {
	type wrapper InMageReplicationDetails
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling InMageReplicationDetails: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling InMageReplicationDetails: %+v", err)
	}
	decoded["instanceType"] = "InMage"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling InMageReplicationDetails: %+v", err)
	}

	return encoded, nil
}

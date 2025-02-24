package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ReplicationProviderSpecificSettings = InMageAzureV2ReplicationDetails{}

type InMageAzureV2ReplicationDetails struct {
	AgentExpiryDate                     *string                                            `json:"agentExpiryDate,omitempty"`
	AgentVersion                        *string                                            `json:"agentVersion,omitempty"`
	AllAvailableOSUpgradeConfigurations *[]OSUpgradeSupportedVersions                      `json:"allAvailableOSUpgradeConfigurations,omitempty"`
	AzureVMDiskDetails                  *[]AzureVMDiskDetails                              `json:"azureVMDiskDetails,omitempty"`
	AzureVMGeneration                   *string                                            `json:"azureVmGeneration,omitempty"`
	CompressedDataRateInMB              *float64                                           `json:"compressedDataRateInMB,omitempty"`
	DataStores                          *[]string                                          `json:"datastores,omitempty"`
	DiscoveryType                       *string                                            `json:"discoveryType,omitempty"`
	DiskResized                         *string                                            `json:"diskResized,omitempty"`
	EnableRdpOnTargetOption             *string                                            `json:"enableRdpOnTargetOption,omitempty"`
	FirmwareType                        *string                                            `json:"firmwareType,omitempty"`
	IPAddress                           *string                                            `json:"ipAddress,omitempty"`
	InfrastructureVMId                  *string                                            `json:"infrastructureVmId,omitempty"`
	IsAdditionalStatsAvailable          *bool                                              `json:"isAdditionalStatsAvailable,omitempty"`
	IsAgentUpdateRequired               *string                                            `json:"isAgentUpdateRequired,omitempty"`
	IsRebootAfterUpdateRequired         *string                                            `json:"isRebootAfterUpdateRequired,omitempty"`
	LastHeartbeat                       *string                                            `json:"lastHeartbeat,omitempty"`
	LastRecoveryPointReceived           *string                                            `json:"lastRecoveryPointReceived,omitempty"`
	LastRpoCalculatedTime               *string                                            `json:"lastRpoCalculatedTime,omitempty"`
	LastUpdateReceivedTime              *string                                            `json:"lastUpdateReceivedTime,omitempty"`
	LicenseType                         *string                                            `json:"licenseType,omitempty"`
	MasterTargetId                      *string                                            `json:"masterTargetId,omitempty"`
	MultiVMGroupId                      *string                                            `json:"multiVmGroupId,omitempty"`
	MultiVMGroupName                    *string                                            `json:"multiVmGroupName,omitempty"`
	MultiVMSyncStatus                   *string                                            `json:"multiVmSyncStatus,omitempty"`
	OsDiskId                            *string                                            `json:"osDiskId,omitempty"`
	OsName                              *string                                            `json:"osName,omitempty"`
	OsType                              *string                                            `json:"osType,omitempty"`
	OsVersion                           *string                                            `json:"osVersion,omitempty"`
	ProcessServerId                     *string                                            `json:"processServerId,omitempty"`
	ProcessServerName                   *string                                            `json:"processServerName,omitempty"`
	ProtectedDisks                      *[]InMageAzureV2ProtectedDiskDetails               `json:"protectedDisks,omitempty"`
	ProtectedManagedDisks               *[]InMageAzureV2ManagedDiskDetails                 `json:"protectedManagedDisks,omitempty"`
	ProtectionStage                     *string                                            `json:"protectionStage,omitempty"`
	RecoveryAvailabilitySetId           *string                                            `json:"recoveryAvailabilitySetId,omitempty"`
	RecoveryAzureLogStorageAccountId    *string                                            `json:"recoveryAzureLogStorageAccountId,omitempty"`
	RecoveryAzureResourceGroupId        *string                                            `json:"recoveryAzureResourceGroupId,omitempty"`
	RecoveryAzureStorageAccount         *string                                            `json:"recoveryAzureStorageAccount,omitempty"`
	RecoveryAzureVMName                 *string                                            `json:"recoveryAzureVMName,omitempty"`
	RecoveryAzureVMSize                 *string                                            `json:"recoveryAzureVMSize,omitempty"`
	ReplicaId                           *string                                            `json:"replicaId,omitempty"`
	ResyncProgressPercentage            *int64                                             `json:"resyncProgressPercentage,omitempty"`
	RpoInSeconds                        *int64                                             `json:"rpoInSeconds,omitempty"`
	SeedManagedDiskTags                 *map[string]string                                 `json:"seedManagedDiskTags,omitempty"`
	SelectedRecoveryAzureNetworkId      *string                                            `json:"selectedRecoveryAzureNetworkId,omitempty"`
	SelectedSourceNicId                 *string                                            `json:"selectedSourceNicId,omitempty"`
	SelectedTfoAzureNetworkId           *string                                            `json:"selectedTfoAzureNetworkId,omitempty"`
	SourceVMCPUCount                    *int64                                             `json:"sourceVmCpuCount,omitempty"`
	SourceVMRamSizeInMB                 *int64                                             `json:"sourceVmRamSizeInMB,omitempty"`
	SqlServerLicenseType                *string                                            `json:"sqlServerLicenseType,omitempty"`
	SupportedOSVersions                 *[]string                                          `json:"supportedOSVersions,omitempty"`
	SwitchProviderBlockingErrorDetails  *[]InMageAzureV2SwitchProviderBlockingErrorDetails `json:"switchProviderBlockingErrorDetails,omitempty"`
	SwitchProviderDetails               *InMageAzureV2SwitchProviderDetails                `json:"switchProviderDetails,omitempty"`
	TargetAvailabilityZone              *string                                            `json:"targetAvailabilityZone,omitempty"`
	TargetManagedDiskTags               *map[string]string                                 `json:"targetManagedDiskTags,omitempty"`
	TargetNicTags                       *map[string]string                                 `json:"targetNicTags,omitempty"`
	TargetProximityPlacementGroupId     *string                                            `json:"targetProximityPlacementGroupId,omitempty"`
	TargetVMId                          *string                                            `json:"targetVmId,omitempty"`
	TargetVMTags                        *map[string]string                                 `json:"targetVmTags,omitempty"`
	TotalDataTransferred                *int64                                             `json:"totalDataTransferred,omitempty"`
	TotalProgressHealth                 *string                                            `json:"totalProgressHealth,omitempty"`
	UncompressedDataRateInMB            *float64                                           `json:"uncompressedDataRateInMB,omitempty"`
	UseManagedDisks                     *string                                            `json:"useManagedDisks,omitempty"`
	VCenterInfrastructureId             *string                                            `json:"vCenterInfrastructureId,omitempty"`
	VMId                                *string                                            `json:"vmId,omitempty"`
	VMNics                              *[]VMNicDetails                                    `json:"vmNics,omitempty"`
	VMProtectionState                   *string                                            `json:"vmProtectionState,omitempty"`
	VMProtectionStateDescription        *string                                            `json:"vmProtectionStateDescription,omitempty"`
	ValidationErrors                    *[]HealthError                                     `json:"validationErrors,omitempty"`
	VhdName                             *string                                            `json:"vhdName,omitempty"`

	// Fields inherited from ReplicationProviderSpecificSettings

	InstanceType string `json:"instanceType"`
}

func (s InMageAzureV2ReplicationDetails) ReplicationProviderSpecificSettings() BaseReplicationProviderSpecificSettingsImpl {
	return BaseReplicationProviderSpecificSettingsImpl{
		InstanceType: s.InstanceType,
	}
}

var _ json.Marshaler = InMageAzureV2ReplicationDetails{}

func (s InMageAzureV2ReplicationDetails) MarshalJSON() ([]byte, error) {
	type wrapper InMageAzureV2ReplicationDetails
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling InMageAzureV2ReplicationDetails: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling InMageAzureV2ReplicationDetails: %+v", err)
	}

	decoded["instanceType"] = "InMageAzureV2"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling InMageAzureV2ReplicationDetails: %+v", err)
	}

	return encoded, nil
}

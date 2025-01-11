package replicationprotecteditems

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/edgezones"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ReplicationProviderSpecificSettings = A2AReplicationDetails{}

type A2AReplicationDetails struct {
	AgentCertificateExpiryDate                  *string                            `json:"agentCertificateExpiryDate,omitempty"`
	AgentExpiryDate                             *string                            `json:"agentExpiryDate,omitempty"`
	AgentVersion                                *string                            `json:"agentVersion,omitempty"`
	AutoProtectionOfDataDisk                    *AutoProtectionOfDataDisk          `json:"autoProtectionOfDataDisk,omitempty"`
	ChurnOptionSelected                         *ChurnOptionSelected               `json:"churnOptionSelected,omitempty"`
	FabricObjectId                              *string                            `json:"fabricObjectId,omitempty"`
	InitialPrimaryExtendedLocation              *edgezones.Model                   `json:"initialPrimaryExtendedLocation,omitempty"`
	InitialPrimaryFabricLocation                *string                            `json:"initialPrimaryFabricLocation,omitempty"`
	InitialPrimaryZone                          *string                            `json:"initialPrimaryZone,omitempty"`
	InitialRecoveryExtendedLocation             *edgezones.Model                   `json:"initialRecoveryExtendedLocation,omitempty"`
	InitialRecoveryFabricLocation               *string                            `json:"initialRecoveryFabricLocation,omitempty"`
	InitialRecoveryZone                         *string                            `json:"initialRecoveryZone,omitempty"`
	IsClusterInfraReady                         *bool                              `json:"isClusterInfraReady,omitempty"`
	IsReplicationAgentCertificateUpdateRequired *bool                              `json:"isReplicationAgentCertificateUpdateRequired,omitempty"`
	IsReplicationAgentUpdateRequired            *bool                              `json:"isReplicationAgentUpdateRequired,omitempty"`
	LastHeartbeat                               *string                            `json:"lastHeartbeat,omitempty"`
	LastRpoCalculatedTime                       *string                            `json:"lastRpoCalculatedTime,omitempty"`
	LifecycleId                                 *string                            `json:"lifecycleId,omitempty"`
	ManagementId                                *string                            `json:"managementId,omitempty"`
	MonitoringJobType                           *string                            `json:"monitoringJobType,omitempty"`
	MonitoringPercentageCompletion              *int64                             `json:"monitoringPercentageCompletion,omitempty"`
	MultiVMGroupCreateOption                    *MultiVMGroupCreateOption          `json:"multiVmGroupCreateOption,omitempty"`
	MultiVMGroupId                              *string                            `json:"multiVmGroupId,omitempty"`
	MultiVMGroupName                            *string                            `json:"multiVmGroupName,omitempty"`
	OsType                                      *string                            `json:"osType,omitempty"`
	PrimaryAvailabilityZone                     *string                            `json:"primaryAvailabilityZone,omitempty"`
	PrimaryExtendedLocation                     *edgezones.Model                   `json:"primaryExtendedLocation,omitempty"`
	PrimaryFabricLocation                       *string                            `json:"primaryFabricLocation,omitempty"`
	ProtectedDisks                              *[]A2AProtectedDiskDetails         `json:"protectedDisks,omitempty"`
	ProtectedManagedDisks                       *[]A2AProtectedManagedDiskDetails  `json:"protectedManagedDisks,omitempty"`
	ProtectionClusterId                         *string                            `json:"protectionClusterId,omitempty"`
	RecoveryAvailabilitySet                     *string                            `json:"recoveryAvailabilitySet,omitempty"`
	RecoveryAvailabilityZone                    *string                            `json:"recoveryAvailabilityZone,omitempty"`
	RecoveryAzureGeneration                     *string                            `json:"recoveryAzureGeneration,omitempty"`
	RecoveryAzureResourceGroupId                *string                            `json:"recoveryAzureResourceGroupId,omitempty"`
	RecoveryAzureVMName                         *string                            `json:"recoveryAzureVMName,omitempty"`
	RecoveryAzureVMSize                         *string                            `json:"recoveryAzureVMSize,omitempty"`
	RecoveryBootDiagStorageAccountId            *string                            `json:"recoveryBootDiagStorageAccountId,omitempty"`
	RecoveryCapacityReservationGroupId          *string                            `json:"recoveryCapacityReservationGroupId,omitempty"`
	RecoveryCloudService                        *string                            `json:"recoveryCloudService,omitempty"`
	RecoveryExtendedLocation                    *edgezones.Model                   `json:"recoveryExtendedLocation,omitempty"`
	RecoveryFabricLocation                      *string                            `json:"recoveryFabricLocation,omitempty"`
	RecoveryFabricObjectId                      *string                            `json:"recoveryFabricObjectId,omitempty"`
	RecoveryProximityPlacementGroupId           *string                            `json:"recoveryProximityPlacementGroupId,omitempty"`
	RecoveryVirtualMachineScaleSetId            *string                            `json:"recoveryVirtualMachineScaleSetId,omitempty"`
	RpoInSeconds                                *int64                             `json:"rpoInSeconds,omitempty"`
	SelectedRecoveryAzureNetworkId              *string                            `json:"selectedRecoveryAzureNetworkId,omitempty"`
	SelectedTfoAzureNetworkId                   *string                            `json:"selectedTfoAzureNetworkId,omitempty"`
	TestFailoverRecoveryFabricObjectId          *string                            `json:"testFailoverRecoveryFabricObjectId,omitempty"`
	TfoAzureVMName                              *string                            `json:"tfoAzureVMName,omitempty"`
	UnprotectedDisks                            *[]A2AUnprotectedDiskDetails       `json:"unprotectedDisks,omitempty"`
	VMEncryptionType                            *VMEncryptionType                  `json:"vmEncryptionType,omitempty"`
	VMNics                                      *[]VMNicDetails                    `json:"vmNics,omitempty"`
	VMProtectionState                           *string                            `json:"vmProtectionState,omitempty"`
	VMProtectionStateDescription                *string                            `json:"vmProtectionStateDescription,omitempty"`
	VMSyncedConfigDetails                       *AzureToAzureVMSyncedConfigDetails `json:"vmSyncedConfigDetails,omitempty"`

	// Fields inherited from ReplicationProviderSpecificSettings

	InstanceType string `json:"instanceType"`
}

func (s A2AReplicationDetails) ReplicationProviderSpecificSettings() BaseReplicationProviderSpecificSettingsImpl {
	return BaseReplicationProviderSpecificSettingsImpl{
		InstanceType: s.InstanceType,
	}
}

var _ json.Marshaler = A2AReplicationDetails{}

func (s A2AReplicationDetails) MarshalJSON() ([]byte, error) {
	type wrapper A2AReplicationDetails
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling A2AReplicationDetails: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling A2AReplicationDetails: %+v", err)
	}

	decoded["instanceType"] = "A2A"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling A2AReplicationDetails: %+v", err)
	}

	return encoded, nil
}

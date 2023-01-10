package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ReplicationProviderSpecificSettings = HyperVReplicaAzureReplicationDetails{}

type HyperVReplicaAzureReplicationDetails struct {
	AzureVMDiskDetails               *[]AzureVMDiskDetails                   `json:"azureVmDiskDetails,omitempty"`
	EnableRdpOnTargetOption          *string                                 `json:"enableRdpOnTargetOption,omitempty"`
	Encryption                       *string                                 `json:"encryption,omitempty"`
	InitialReplicationDetails        *InitialReplicationDetails              `json:"initialReplicationDetails,omitempty"`
	LastRecoveryPointReceived        *string                                 `json:"lastRecoveryPointReceived,omitempty"`
	LastReplicatedTime               *string                                 `json:"lastReplicatedTime,omitempty"`
	LastRpoCalculatedTime            *string                                 `json:"lastRpoCalculatedTime,omitempty"`
	LicenseType                      *string                                 `json:"licenseType,omitempty"`
	OSDetails                        *OSDetails                              `json:"oSDetails,omitempty"`
	ProtectedManagedDisks            *[]HyperVReplicaAzureManagedDiskDetails `json:"protectedManagedDisks,omitempty"`
	RecoveryAvailabilitySetId        *string                                 `json:"recoveryAvailabilitySetId,omitempty"`
	RecoveryAzureLogStorageAccountId *string                                 `json:"recoveryAzureLogStorageAccountId,omitempty"`
	RecoveryAzureResourceGroupId     *string                                 `json:"recoveryAzureResourceGroupId,omitempty"`
	RecoveryAzureStorageAccount      *string                                 `json:"recoveryAzureStorageAccount,omitempty"`
	RecoveryAzureVMName              *string                                 `json:"recoveryAzureVmName,omitempty"`
	RecoveryAzureVMSize              *string                                 `json:"recoveryAzureVMSize,omitempty"`
	RpoInSeconds                     *int64                                  `json:"rpoInSeconds,omitempty"`
	SeedManagedDiskTags              *map[string]string                      `json:"seedManagedDiskTags,omitempty"`
	SelectedRecoveryAzureNetworkId   *string                                 `json:"selectedRecoveryAzureNetworkId,omitempty"`
	SelectedSourceNicId              *string                                 `json:"selectedSourceNicId,omitempty"`
	SourceVMCPUCount                 *int64                                  `json:"sourceVmCpuCount,omitempty"`
	SourceVMRamSizeInMB              *int64                                  `json:"sourceVmRamSizeInMB,omitempty"`
	SqlServerLicenseType             *string                                 `json:"sqlServerLicenseType,omitempty"`
	TargetAvailabilityZone           *string                                 `json:"targetAvailabilityZone,omitempty"`
	TargetManagedDiskTags            *map[string]string                      `json:"targetManagedDiskTags,omitempty"`
	TargetNicTags                    *map[string]string                      `json:"targetNicTags,omitempty"`
	TargetProximityPlacementGroupId  *string                                 `json:"targetProximityPlacementGroupId,omitempty"`
	TargetVMTags                     *map[string]string                      `json:"targetVmTags,omitempty"`
	UseManagedDisks                  *string                                 `json:"useManagedDisks,omitempty"`
	VMId                             *string                                 `json:"vmId,omitempty"`
	VMNics                           *[]VMNicDetails                         `json:"vmNics,omitempty"`
	VMProtectionState                *string                                 `json:"vmProtectionState,omitempty"`
	VMProtectionStateDescription     *string                                 `json:"vmProtectionStateDescription,omitempty"`

	// Fields inherited from ReplicationProviderSpecificSettings
}

var _ json.Marshaler = HyperVReplicaAzureReplicationDetails{}

func (s HyperVReplicaAzureReplicationDetails) MarshalJSON() ([]byte, error) {
	type wrapper HyperVReplicaAzureReplicationDetails
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling HyperVReplicaAzureReplicationDetails: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling HyperVReplicaAzureReplicationDetails: %+v", err)
	}
	decoded["instanceType"] = "HyperVReplicaAzure"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling HyperVReplicaAzureReplicationDetails: %+v", err)
	}

	return encoded, nil
}

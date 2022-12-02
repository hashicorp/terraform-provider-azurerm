package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ReplicationProviderSpecificSettings = HyperVReplicaAzureReplicationDetails{}

type HyperVReplicaAzureReplicationDetails struct {
	AzureVmDiskDetails               *[]AzureVmDiskDetails                   `json:"azureVmDiskDetails,omitempty"`
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
	RecoveryAzureVMSize              *string                                 `json:"recoveryAzureVMSize,omitempty"`
	RecoveryAzureVmName              *string                                 `json:"recoveryAzureVmName,omitempty"`
	RpoInSeconds                     *int64                                  `json:"rpoInSeconds,omitempty"`
	SeedManagedDiskTags              *map[string]string                      `json:"seedManagedDiskTags,omitempty"`
	SelectedRecoveryAzureNetworkId   *string                                 `json:"selectedRecoveryAzureNetworkId,omitempty"`
	SelectedSourceNicId              *string                                 `json:"selectedSourceNicId,omitempty"`
	SourceVmCpuCount                 *int64                                  `json:"sourceVmCpuCount,omitempty"`
	SourceVmRamSizeInMB              *int64                                  `json:"sourceVmRamSizeInMB,omitempty"`
	SqlServerLicenseType             *string                                 `json:"sqlServerLicenseType,omitempty"`
	TargetAvailabilityZone           *string                                 `json:"targetAvailabilityZone,omitempty"`
	TargetManagedDiskTags            *map[string]string                      `json:"targetManagedDiskTags,omitempty"`
	TargetNicTags                    *map[string]string                      `json:"targetNicTags,omitempty"`
	TargetProximityPlacementGroupId  *string                                 `json:"targetProximityPlacementGroupId,omitempty"`
	TargetVmTags                     *map[string]string                      `json:"targetVmTags,omitempty"`
	UseManagedDisks                  *string                                 `json:"useManagedDisks,omitempty"`
	VmId                             *string                                 `json:"vmId,omitempty"`
	VmNics                           *[]VMNicDetails                         `json:"vmNics,omitempty"`
	VmProtectionState                *string                                 `json:"vmProtectionState,omitempty"`
	VmProtectionStateDescription     *string                                 `json:"vmProtectionStateDescription,omitempty"`

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

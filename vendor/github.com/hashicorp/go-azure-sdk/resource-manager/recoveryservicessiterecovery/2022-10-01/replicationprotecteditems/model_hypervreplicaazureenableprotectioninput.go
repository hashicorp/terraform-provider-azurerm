package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ EnableProtectionProviderSpecificInput = HyperVReplicaAzureEnableProtectionInput{}

type HyperVReplicaAzureEnableProtectionInput struct {
	DiskEncryptionSetId             *string                               `json:"diskEncryptionSetId,omitempty"`
	DiskType                        *DiskAccountType                      `json:"diskType,omitempty"`
	DisksToInclude                  *[]string                             `json:"disksToInclude,omitempty"`
	DisksToIncludeForManagedDisks   *[]HyperVReplicaAzureDiskInputDetails `json:"disksToIncludeForManagedDisks,omitempty"`
	EnableRdpOnTargetOption         *string                               `json:"enableRdpOnTargetOption,omitempty"`
	HvHostVMId                      *string                               `json:"hvHostVmId,omitempty"`
	LicenseType                     *LicenseType                          `json:"licenseType,omitempty"`
	LogStorageAccountId             *string                               `json:"logStorageAccountId,omitempty"`
	OsType                          *string                               `json:"osType,omitempty"`
	SeedManagedDiskTags             *map[string]string                    `json:"seedManagedDiskTags,omitempty"`
	SqlServerLicenseType            *SqlServerLicenseType                 `json:"sqlServerLicenseType,omitempty"`
	TargetAvailabilitySetId         *string                               `json:"targetAvailabilitySetId,omitempty"`
	TargetAvailabilityZone          *string                               `json:"targetAvailabilityZone,omitempty"`
	TargetAzureNetworkId            *string                               `json:"targetAzureNetworkId,omitempty"`
	TargetAzureSubnetId             *string                               `json:"targetAzureSubnetId,omitempty"`
	TargetAzureV1ResourceGroupId    *string                               `json:"targetAzureV1ResourceGroupId,omitempty"`
	TargetAzureV2ResourceGroupId    *string                               `json:"targetAzureV2ResourceGroupId,omitempty"`
	TargetAzureVMName               *string                               `json:"targetAzureVmName,omitempty"`
	TargetManagedDiskTags           *map[string]string                    `json:"targetManagedDiskTags,omitempty"`
	TargetNicTags                   *map[string]string                    `json:"targetNicTags,omitempty"`
	TargetProximityPlacementGroupId *string                               `json:"targetProximityPlacementGroupId,omitempty"`
	TargetStorageAccountId          *string                               `json:"targetStorageAccountId,omitempty"`
	TargetVMSize                    *string                               `json:"targetVmSize,omitempty"`
	TargetVMTags                    *map[string]string                    `json:"targetVmTags,omitempty"`
	UseManagedDisks                 *string                               `json:"useManagedDisks,omitempty"`
	UseManagedDisksForReplication   *string                               `json:"useManagedDisksForReplication,omitempty"`
	VhdId                           *string                               `json:"vhdId,omitempty"`
	VirtualMachineName              *string                               `json:"vmName,omitempty"`

	// Fields inherited from EnableProtectionProviderSpecificInput
}

var _ json.Marshaler = HyperVReplicaAzureEnableProtectionInput{}

func (s HyperVReplicaAzureEnableProtectionInput) MarshalJSON() ([]byte, error) {
	type wrapper HyperVReplicaAzureEnableProtectionInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling HyperVReplicaAzureEnableProtectionInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling HyperVReplicaAzureEnableProtectionInput: %+v", err)
	}
	decoded["instanceType"] = "HyperVReplicaAzure"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling HyperVReplicaAzureEnableProtectionInput: %+v", err)
	}

	return encoded, nil
}

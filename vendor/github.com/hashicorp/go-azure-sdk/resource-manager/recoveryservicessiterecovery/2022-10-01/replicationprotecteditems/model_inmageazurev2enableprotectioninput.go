package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ EnableProtectionProviderSpecificInput = InMageAzureV2EnableProtectionInput{}

type InMageAzureV2EnableProtectionInput struct {
	DiskEncryptionSetId             *string                          `json:"diskEncryptionSetId,omitempty"`
	DiskType                        *DiskAccountType                 `json:"diskType,omitempty"`
	DisksToInclude                  *[]InMageAzureV2DiskInputDetails `json:"disksToInclude,omitempty"`
	EnableRdpOnTargetOption         *string                          `json:"enableRdpOnTargetOption,omitempty"`
	LicenseType                     *LicenseType                     `json:"licenseType,omitempty"`
	LogStorageAccountId             *string                          `json:"logStorageAccountId,omitempty"`
	MasterTargetId                  *string                          `json:"masterTargetId,omitempty"`
	MultiVMGroupId                  *string                          `json:"multiVmGroupId,omitempty"`
	MultiVMGroupName                *string                          `json:"multiVmGroupName,omitempty"`
	ProcessServerId                 *string                          `json:"processServerId,omitempty"`
	RunAsAccountId                  *string                          `json:"runAsAccountId,omitempty"`
	SeedManagedDiskTags             *map[string]string               `json:"seedManagedDiskTags,omitempty"`
	SqlServerLicenseType            *SqlServerLicenseType            `json:"sqlServerLicenseType,omitempty"`
	StorageAccountId                *string                          `json:"storageAccountId,omitempty"`
	TargetAvailabilitySetId         *string                          `json:"targetAvailabilitySetId,omitempty"`
	TargetAvailabilityZone          *string                          `json:"targetAvailabilityZone,omitempty"`
	TargetAzureNetworkId            *string                          `json:"targetAzureNetworkId,omitempty"`
	TargetAzureSubnetId             *string                          `json:"targetAzureSubnetId,omitempty"`
	TargetAzureV1ResourceGroupId    *string                          `json:"targetAzureV1ResourceGroupId,omitempty"`
	TargetAzureV2ResourceGroupId    *string                          `json:"targetAzureV2ResourceGroupId,omitempty"`
	TargetAzureVMName               *string                          `json:"targetAzureVmName,omitempty"`
	TargetManagedDiskTags           *map[string]string               `json:"targetManagedDiskTags,omitempty"`
	TargetNicTags                   *map[string]string               `json:"targetNicTags,omitempty"`
	TargetProximityPlacementGroupId *string                          `json:"targetProximityPlacementGroupId,omitempty"`
	TargetVMSize                    *string                          `json:"targetVmSize,omitempty"`
	TargetVMTags                    *map[string]string               `json:"targetVmTags,omitempty"`

	// Fields inherited from EnableProtectionProviderSpecificInput
}

var _ json.Marshaler = InMageAzureV2EnableProtectionInput{}

func (s InMageAzureV2EnableProtectionInput) MarshalJSON() ([]byte, error) {
	type wrapper InMageAzureV2EnableProtectionInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling InMageAzureV2EnableProtectionInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling InMageAzureV2EnableProtectionInput: %+v", err)
	}
	decoded["instanceType"] = "InMageAzureV2"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling InMageAzureV2EnableProtectionInput: %+v", err)
	}

	return encoded, nil
}

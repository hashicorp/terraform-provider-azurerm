package replicationprotecteditems

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/edgezones"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ EnableProtectionProviderSpecificInput = A2AEnableProtectionInput{}

type A2AEnableProtectionInput struct {
	DiskEncryptionInfo                 *DiskEncryptionInfo             `json:"diskEncryptionInfo,omitempty"`
	FabricObjectId                     string                          `json:"fabricObjectId"`
	MultiVMGroupId                     *string                         `json:"multiVmGroupId,omitempty"`
	MultiVMGroupName                   *string                         `json:"multiVmGroupName,omitempty"`
	RecoveryAvailabilitySetId          *string                         `json:"recoveryAvailabilitySetId,omitempty"`
	RecoveryAvailabilityZone           *string                         `json:"recoveryAvailabilityZone,omitempty"`
	RecoveryAzureNetworkId             *string                         `json:"recoveryAzureNetworkId,omitempty"`
	RecoveryBootDiagStorageAccountId   *string                         `json:"recoveryBootDiagStorageAccountId,omitempty"`
	RecoveryCapacityReservationGroupId *string                         `json:"recoveryCapacityReservationGroupId,omitempty"`
	RecoveryCloudServiceId             *string                         `json:"recoveryCloudServiceId,omitempty"`
	RecoveryContainerId                *string                         `json:"recoveryContainerId,omitempty"`
	RecoveryExtendedLocation           *edgezones.Model                `json:"recoveryExtendedLocation,omitempty"`
	RecoveryProximityPlacementGroupId  *string                         `json:"recoveryProximityPlacementGroupId,omitempty"`
	RecoveryResourceGroupId            *string                         `json:"recoveryResourceGroupId,omitempty"`
	RecoverySubnetName                 *string                         `json:"recoverySubnetName,omitempty"`
	RecoveryVirtualMachineScaleSetId   *string                         `json:"recoveryVirtualMachineScaleSetId,omitempty"`
	VMDisks                            *[]A2AVMDiskInputDetails        `json:"vmDisks,omitempty"`
	VMManagedDisks                     *[]A2AVMManagedDiskInputDetails `json:"vmManagedDisks,omitempty"`

	// Fields inherited from EnableProtectionProviderSpecificInput
}

var _ json.Marshaler = A2AEnableProtectionInput{}

func (s A2AEnableProtectionInput) MarshalJSON() ([]byte, error) {
	type wrapper A2AEnableProtectionInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling A2AEnableProtectionInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling A2AEnableProtectionInput: %+v", err)
	}
	decoded["instanceType"] = "A2A"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling A2AEnableProtectionInput: %+v", err)
	}

	return encoded, nil
}

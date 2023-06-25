package replicationprotectioncontainers

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ SwitchProtectionProviderSpecificInput = A2ASwitchProtectionInput{}

type A2ASwitchProtectionInput struct {
	DiskEncryptionInfo                 *DiskEncryptionInfo             `json:"diskEncryptionInfo,omitempty"`
	PolicyId                           *string                         `json:"policyId,omitempty"`
	RecoveryAvailabilitySetId          *string                         `json:"recoveryAvailabilitySetId,omitempty"`
	RecoveryAvailabilityZone           *string                         `json:"recoveryAvailabilityZone,omitempty"`
	RecoveryBootDiagStorageAccountId   *string                         `json:"recoveryBootDiagStorageAccountId,omitempty"`
	RecoveryCapacityReservationGroupId *string                         `json:"recoveryCapacityReservationGroupId,omitempty"`
	RecoveryCloudServiceId             *string                         `json:"recoveryCloudServiceId,omitempty"`
	RecoveryContainerId                *string                         `json:"recoveryContainerId,omitempty"`
	RecoveryProximityPlacementGroupId  *string                         `json:"recoveryProximityPlacementGroupId,omitempty"`
	RecoveryResourceGroupId            *string                         `json:"recoveryResourceGroupId,omitempty"`
	RecoveryVirtualMachineScaleSetId   *string                         `json:"recoveryVirtualMachineScaleSetId,omitempty"`
	VMDisks                            *[]A2AVMDiskInputDetails        `json:"vmDisks,omitempty"`
	VMManagedDisks                     *[]A2AVMManagedDiskInputDetails `json:"vmManagedDisks,omitempty"`

	// Fields inherited from SwitchProtectionProviderSpecificInput
}

var _ json.Marshaler = A2ASwitchProtectionInput{}

func (s A2ASwitchProtectionInput) MarshalJSON() ([]byte, error) {
	type wrapper A2ASwitchProtectionInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling A2ASwitchProtectionInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling A2ASwitchProtectionInput: %+v", err)
	}
	decoded["instanceType"] = "A2A"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling A2ASwitchProtectionInput: %+v", err)
	}

	return encoded, nil
}

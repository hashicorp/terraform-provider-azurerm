package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ UpdateReplicationProtectedItemProviderInput = A2AUpdateReplicationProtectedItemInput{}

type A2AUpdateReplicationProtectedItemInput struct {
	DiskEncryptionInfo                 *DiskEncryptionInfo              `json:"diskEncryptionInfo,omitempty"`
	ManagedDiskUpdateDetails           *[]A2AVMManagedDiskUpdateDetails `json:"managedDiskUpdateDetails,omitempty"`
	RecoveryBootDiagStorageAccountId   *string                          `json:"recoveryBootDiagStorageAccountId,omitempty"`
	RecoveryCapacityReservationGroupId *string                          `json:"recoveryCapacityReservationGroupId,omitempty"`
	RecoveryCloudServiceId             *string                          `json:"recoveryCloudServiceId,omitempty"`
	RecoveryProximityPlacementGroupId  *string                          `json:"recoveryProximityPlacementGroupId,omitempty"`
	RecoveryResourceGroupId            *string                          `json:"recoveryResourceGroupId,omitempty"`
	RecoveryVirtualMachineScaleSetId   *string                          `json:"recoveryVirtualMachineScaleSetId,omitempty"`
	TfoAzureVMName                     *string                          `json:"tfoAzureVMName,omitempty"`

	// Fields inherited from UpdateReplicationProtectedItemProviderInput
}

var _ json.Marshaler = A2AUpdateReplicationProtectedItemInput{}

func (s A2AUpdateReplicationProtectedItemInput) MarshalJSON() ([]byte, error) {
	type wrapper A2AUpdateReplicationProtectedItemInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling A2AUpdateReplicationProtectedItemInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling A2AUpdateReplicationProtectedItemInput: %+v", err)
	}
	decoded["instanceType"] = "A2A"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling A2AUpdateReplicationProtectedItemInput: %+v", err)
	}

	return encoded, nil
}

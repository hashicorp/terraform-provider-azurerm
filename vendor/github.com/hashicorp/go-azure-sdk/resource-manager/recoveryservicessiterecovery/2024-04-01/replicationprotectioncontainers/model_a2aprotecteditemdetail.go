package replicationprotectioncontainers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type A2AProtectedItemDetail struct {
	DiskEncryptionInfo                 *DiskEncryptionInfo             `json:"diskEncryptionInfo,omitempty"`
	RecoveryAvailabilitySetId          *string                         `json:"recoveryAvailabilitySetId,omitempty"`
	RecoveryAvailabilityZone           *string                         `json:"recoveryAvailabilityZone,omitempty"`
	RecoveryBootDiagStorageAccountId   *string                         `json:"recoveryBootDiagStorageAccountId,omitempty"`
	RecoveryCapacityReservationGroupId *string                         `json:"recoveryCapacityReservationGroupId,omitempty"`
	RecoveryProximityPlacementGroupId  *string                         `json:"recoveryProximityPlacementGroupId,omitempty"`
	RecoveryResourceGroupId            *string                         `json:"recoveryResourceGroupId,omitempty"`
	RecoveryVirtualMachineScaleSetId   *string                         `json:"recoveryVirtualMachineScaleSetId,omitempty"`
	ReplicationProtectedItemName       *string                         `json:"replicationProtectedItemName,omitempty"`
	VMManagedDisks                     *[]A2AVMManagedDiskInputDetails `json:"vmManagedDisks,omitempty"`
}

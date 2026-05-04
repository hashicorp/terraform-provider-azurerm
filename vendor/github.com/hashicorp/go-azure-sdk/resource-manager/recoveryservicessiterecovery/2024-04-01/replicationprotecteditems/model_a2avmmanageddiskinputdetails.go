package replicationprotecteditems

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type A2AVMManagedDiskInputDetails struct {
	DiskEncryptionInfo                  *DiskEncryptionInfo `json:"diskEncryptionInfo,omitempty"`
	DiskId                              string              `json:"diskId"`
	PrimaryStagingAzureStorageAccountId string              `json:"primaryStagingAzureStorageAccountId"`
	RecoveryDiskEncryptionSetId         *string             `json:"recoveryDiskEncryptionSetId,omitempty"`
	RecoveryReplicaDiskAccountType      *string             `json:"recoveryReplicaDiskAccountType,omitempty"`
	RecoveryResourceGroupId             string              `json:"recoveryResourceGroupId"`
	RecoveryTargetDiskAccountType       *string             `json:"recoveryTargetDiskAccountType,omitempty"`
}

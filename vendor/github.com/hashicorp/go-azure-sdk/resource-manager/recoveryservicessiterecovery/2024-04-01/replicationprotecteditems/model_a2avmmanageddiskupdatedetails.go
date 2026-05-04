package replicationprotecteditems

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type A2AVMManagedDiskUpdateDetails struct {
	DiskEncryptionInfo             *DiskEncryptionInfo `json:"diskEncryptionInfo,omitempty"`
	DiskId                         *string             `json:"diskId,omitempty"`
	FailoverDiskName               *string             `json:"failoverDiskName,omitempty"`
	RecoveryReplicaDiskAccountType *string             `json:"recoveryReplicaDiskAccountType,omitempty"`
	RecoveryTargetDiskAccountType  *string             `json:"recoveryTargetDiskAccountType,omitempty"`
	TfoDiskName                    *string             `json:"tfoDiskName,omitempty"`
}

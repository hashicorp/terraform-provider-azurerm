package replicationprotecteditems

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InMageRcmProtectedDiskDetails struct {
	CapacityInBytes               *int64                `json:"capacityInBytes,omitempty"`
	DataPendingAtSourceAgentInMB  *float64              `json:"dataPendingAtSourceAgentInMB,omitempty"`
	DataPendingInLogDataStoreInMB *float64              `json:"dataPendingInLogDataStoreInMB,omitempty"`
	DiskEncryptionSetId           *string               `json:"diskEncryptionSetId,omitempty"`
	DiskId                        *string               `json:"diskId,omitempty"`
	DiskName                      *string               `json:"diskName,omitempty"`
	DiskType                      *DiskAccountType      `json:"diskType,omitempty"`
	IrDetails                     *InMageRcmSyncDetails `json:"irDetails,omitempty"`
	IsInitialReplicationComplete  *string               `json:"isInitialReplicationComplete,omitempty"`
	IsOSDisk                      *string               `json:"isOSDisk,omitempty"`
	LogStorageAccountId           *string               `json:"logStorageAccountId,omitempty"`
	ResyncDetails                 *InMageRcmSyncDetails `json:"resyncDetails,omitempty"`
	SeedBlobUri                   *string               `json:"seedBlobUri,omitempty"`
	SeedManagedDiskId             *string               `json:"seedManagedDiskId,omitempty"`
	TargetManagedDiskId           *string               `json:"targetManagedDiskId,omitempty"`
}

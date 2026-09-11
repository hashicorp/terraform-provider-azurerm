package managedclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedClusterStorageProfile struct {
	BlobCSIDriver      *ManagedClusterStorageProfileBlobCSIDriver      `json:"blobCSIDriver,omitempty"`
	DiskCSIDriver      *ManagedClusterStorageProfileDiskCSIDriver      `json:"diskCSIDriver,omitempty"`
	FileCSIDriver      *ManagedClusterStorageProfileFileCSIDriver      `json:"fileCSIDriver,omitempty"`
	SnapshotController *ManagedClusterStorageProfileSnapshotController `json:"snapshotController,omitempty"`
}

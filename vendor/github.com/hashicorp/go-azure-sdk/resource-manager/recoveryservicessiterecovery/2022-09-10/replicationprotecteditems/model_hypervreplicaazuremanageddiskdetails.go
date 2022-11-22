package replicationprotecteditems

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HyperVReplicaAzureManagedDiskDetails struct {
	DiskEncryptionSetId *string `json:"diskEncryptionSetId,omitempty"`
	DiskId              *string `json:"diskId,omitempty"`
	ReplicaDiskType     *string `json:"replicaDiskType,omitempty"`
	SeedManagedDiskId   *string `json:"seedManagedDiskId,omitempty"`
}

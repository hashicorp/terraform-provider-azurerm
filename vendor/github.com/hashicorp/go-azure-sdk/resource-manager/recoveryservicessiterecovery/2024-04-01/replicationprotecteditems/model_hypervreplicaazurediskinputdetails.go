package replicationprotecteditems

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HyperVReplicaAzureDiskInputDetails struct {
	DiskEncryptionSetId *string          `json:"diskEncryptionSetId,omitempty"`
	DiskId              *string          `json:"diskId,omitempty"`
	DiskType            *DiskAccountType `json:"diskType,omitempty"`
	LogStorageAccountId *string          `json:"logStorageAccountId,omitempty"`
	SectorSizeInBytes   *int64           `json:"sectorSizeInBytes,omitempty"`
}

package replicationprotecteditems

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InMageRcmDiskInput struct {
	DiskEncryptionSetId *string         `json:"diskEncryptionSetId,omitempty"`
	DiskId              string          `json:"diskId"`
	DiskType            DiskAccountType `json:"diskType"`
	LogStorageAccountId string          `json:"logStorageAccountId"`
	SectorSizeInBytes   *int64          `json:"sectorSizeInBytes,omitempty"`
}

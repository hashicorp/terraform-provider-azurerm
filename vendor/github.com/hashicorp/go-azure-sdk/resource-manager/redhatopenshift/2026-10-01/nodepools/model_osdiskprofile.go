package nodepools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OsDiskProfile struct {
	DiskStorageAccountType *DiskStorageAccountType `json:"diskStorageAccountType,omitempty"`
	DiskType               *OsDiskType             `json:"diskType,omitempty"`
	EncryptionSetId        *string                 `json:"encryptionSetId,omitempty"`
	SizeGiB                *int64                  `json:"sizeGiB,omitempty"`
}

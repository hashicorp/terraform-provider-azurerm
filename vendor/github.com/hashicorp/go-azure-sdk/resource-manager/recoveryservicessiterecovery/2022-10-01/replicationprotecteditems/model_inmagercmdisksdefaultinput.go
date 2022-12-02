package replicationprotecteditems

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InMageRcmDisksDefaultInput struct {
	DiskEncryptionSetId *string         `json:"diskEncryptionSetId,omitempty"`
	DiskType            DiskAccountType `json:"diskType"`
	LogStorageAccountId string          `json:"logStorageAccountId"`
}

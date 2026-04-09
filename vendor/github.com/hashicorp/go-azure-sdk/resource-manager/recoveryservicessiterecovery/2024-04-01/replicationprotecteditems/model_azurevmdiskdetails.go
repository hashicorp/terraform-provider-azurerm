package replicationprotecteditems

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureVMDiskDetails struct {
	CustomTargetDiskName *string `json:"customTargetDiskName,omitempty"`
	DiskEncryptionSetId  *string `json:"diskEncryptionSetId,omitempty"`
	DiskId               *string `json:"diskId,omitempty"`
	LunId                *string `json:"lunId,omitempty"`
	MaxSizeMB            *string `json:"maxSizeMB,omitempty"`
	TargetDiskLocation   *string `json:"targetDiskLocation,omitempty"`
	TargetDiskName       *string `json:"targetDiskName,omitempty"`
	VhdId                *string `json:"vhdId,omitempty"`
	VhdName              *string `json:"vhdName,omitempty"`
	VhdType              *string `json:"vhdType,omitempty"`
}

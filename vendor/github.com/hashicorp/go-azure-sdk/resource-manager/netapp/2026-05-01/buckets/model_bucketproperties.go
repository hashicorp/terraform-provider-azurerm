package buckets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BucketProperties struct {
	AkvDetails        *AzureKeyVaultDetails    `json:"akvDetails,omitempty"`
	FileSystemUser    *FileSystemUser          `json:"fileSystemUser,omitempty"`
	Path              *string                  `json:"path,omitempty"`
	Permissions       *BucketPermissions       `json:"permissions,omitempty"`
	ProvisioningState *NetAppProvisioningState `json:"provisioningState,omitempty"`
	Server            *BucketServerProperties  `json:"server,omitempty"`
	Status            *CredentialsStatus       `json:"status,omitempty"`
}

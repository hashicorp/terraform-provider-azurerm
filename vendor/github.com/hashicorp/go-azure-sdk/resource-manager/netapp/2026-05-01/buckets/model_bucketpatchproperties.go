package buckets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BucketPatchProperties struct {
	AkvDetails        *AzureKeyVaultDetails        `json:"akvDetails,omitempty"`
	FileSystemUser    *FileSystemUser              `json:"fileSystemUser,omitempty"`
	Permissions       *BucketPatchPermissions      `json:"permissions,omitempty"`
	ProvisioningState *NetAppProvisioningState     `json:"provisioningState,omitempty"`
	Server            *BucketServerPatchProperties `json:"server,omitempty"`
}

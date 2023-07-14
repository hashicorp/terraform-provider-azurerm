package pool

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MountConfiguration struct {
	AzureBlobFileSystemConfiguration *AzureBlobFileSystemConfiguration `json:"azureBlobFileSystemConfiguration,omitempty"`
	AzureFileShareConfiguration      *AzureFileShareConfiguration      `json:"azureFileShareConfiguration,omitempty"`
	CifsMountConfiguration           *CIFSMountConfiguration           `json:"cifsMountConfiguration,omitempty"`
	NfsMountConfiguration            *NFSMountConfiguration            `json:"nfsMountConfiguration,omitempty"`
}

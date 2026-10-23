package cloudhsmclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackupRestoreRequestBaseProperties struct {
	AzureStorageBlobContainerUri string  `json:"azureStorageBlobContainerUri"`
	Token                        *string `json:"token,omitempty"`
}

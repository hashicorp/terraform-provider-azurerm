package cloudhsmclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RestoreRequestProperties struct {
	AzureStorageBlobContainerUri string  `json:"azureStorageBlobContainerUri"`
	BackupId                     string  `json:"backupId"`
	Token                        *string `json:"token,omitempty"`
}

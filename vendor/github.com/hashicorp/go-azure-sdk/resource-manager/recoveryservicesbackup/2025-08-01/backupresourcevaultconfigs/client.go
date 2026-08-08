package backupresourcevaultconfigs

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackupResourceVaultConfigsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewBackupResourceVaultConfigsClientWithBaseURI(endpoint string) BackupResourceVaultConfigsClient {
	return BackupResourceVaultConfigsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

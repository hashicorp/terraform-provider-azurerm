package backupprotecteditems

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackupProtectedItemsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewBackupProtectedItemsClientWithBaseURI(endpoint string) BackupProtectedItemsClient {
	return BackupProtectedItemsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

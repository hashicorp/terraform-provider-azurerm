package backupresourcestorageconfigsnoncrr

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackupResourceStorageConfigsNonCRRClient struct {
	Client  autorest.Client
	baseUri string
}

func NewBackupResourceStorageConfigsNonCRRClientWithBaseURI(endpoint string) BackupResourceStorageConfigsNonCRRClient {
	return BackupResourceStorageConfigsNonCRRClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

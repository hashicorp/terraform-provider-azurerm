package storageaccounts

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageAccountsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewStorageAccountsClientWithBaseURI(endpoint string) StorageAccountsClient {
	return StorageAccountsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

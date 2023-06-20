package storageinsights

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageInsightsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewStorageInsightsClientWithBaseURI(endpoint string) StorageInsightsClient {
	return StorageInsightsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

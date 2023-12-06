package defenderforstorage

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DefenderForStorageClient struct {
	Client  autorest.Client
	baseUri string
}

func NewDefenderForStorageClientWithBaseURI(endpoint string) DefenderForStorageClient {
	return DefenderForStorageClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

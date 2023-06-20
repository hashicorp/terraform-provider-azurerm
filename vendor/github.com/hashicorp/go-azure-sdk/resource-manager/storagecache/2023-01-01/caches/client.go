package caches

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CachesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewCachesClientWithBaseURI(endpoint string) CachesClient {
	return CachesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

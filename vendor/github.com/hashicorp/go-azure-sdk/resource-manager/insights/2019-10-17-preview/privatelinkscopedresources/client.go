package privatelinkscopedresources

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateLinkScopedResourcesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewPrivateLinkScopedResourcesClientWithBaseURI(endpoint string) PrivateLinkScopedResourcesClient {
	return PrivateLinkScopedResourcesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

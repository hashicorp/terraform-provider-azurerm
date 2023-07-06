package privatelinkscopes

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateLinkScopesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewPrivateLinkScopesClientWithBaseURI(endpoint string) PrivateLinkScopesClient {
	return PrivateLinkScopesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

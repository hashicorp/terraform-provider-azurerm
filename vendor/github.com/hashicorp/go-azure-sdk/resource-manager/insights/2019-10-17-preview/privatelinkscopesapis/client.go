package privatelinkscopesapis

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateLinkScopesAPIsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewPrivateLinkScopesAPIsClientWithBaseURI(endpoint string) PrivateLinkScopesAPIsClient {
	return PrivateLinkScopesAPIsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

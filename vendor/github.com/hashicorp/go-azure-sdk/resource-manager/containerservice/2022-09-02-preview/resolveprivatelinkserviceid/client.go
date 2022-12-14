package resolveprivatelinkserviceid

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResolvePrivateLinkServiceIdClient struct {
	Client  autorest.Client
	baseUri string
}

func NewResolvePrivateLinkServiceIdClientWithBaseURI(endpoint string) ResolvePrivateLinkServiceIdClient {
	return ResolvePrivateLinkServiceIdClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

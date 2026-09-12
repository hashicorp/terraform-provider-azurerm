package resourceguardproxies

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceGuardProxiesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewResourceGuardProxiesClientWithBaseURI(endpoint string) ResourceGuardProxiesClient {
	return ResourceGuardProxiesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

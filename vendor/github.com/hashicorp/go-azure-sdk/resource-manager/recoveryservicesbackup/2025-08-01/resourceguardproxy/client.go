package resourceguardproxy

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceGuardProxyClient struct {
	Client  autorest.Client
	baseUri string
}

func NewResourceGuardProxyClientWithBaseURI(endpoint string) ResourceGuardProxyClient {
	return ResourceGuardProxyClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

package privatelinkresources

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateLinkResourcesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewPrivateLinkResourcesClientWithBaseURI(endpoint string) PrivateLinkResourcesClient {
	return PrivateLinkResourcesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

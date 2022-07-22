package sharedprivatelinkresources

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SharedPrivateLinkResourcesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewSharedPrivateLinkResourcesClientWithBaseURI(endpoint string) SharedPrivateLinkResourcesClient {
	return SharedPrivateLinkResourcesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

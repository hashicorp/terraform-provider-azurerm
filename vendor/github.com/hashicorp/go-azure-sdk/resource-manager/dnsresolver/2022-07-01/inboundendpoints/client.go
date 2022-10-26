package inboundendpoints

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InboundEndpointsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewInboundEndpointsClientWithBaseURI(endpoint string) InboundEndpointsClient {
	return InboundEndpointsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

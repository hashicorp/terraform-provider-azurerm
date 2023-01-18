package streamingendpoint

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StreamingEndpointClient struct {
	Client  autorest.Client
	baseUri string
}

func NewStreamingEndpointClientWithBaseURI(endpoint string) StreamingEndpointClient {
	return StreamingEndpointClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

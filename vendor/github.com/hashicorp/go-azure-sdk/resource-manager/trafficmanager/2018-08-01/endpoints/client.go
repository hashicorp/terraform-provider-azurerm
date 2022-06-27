package endpoints

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EndpointsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewEndpointsClientWithBaseURI(endpoint string) EndpointsClient {
	return EndpointsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

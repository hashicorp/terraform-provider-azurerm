package providers

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProvidersClient struct {
	Client  autorest.Client
	baseUri string
}

func NewProvidersClientWithBaseURI(endpoint string) ProvidersClient {
	return ProvidersClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

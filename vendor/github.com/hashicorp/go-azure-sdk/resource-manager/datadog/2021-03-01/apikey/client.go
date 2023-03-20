package apikey

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiKeyClient struct {
	Client  autorest.Client
	baseUri string
}

func NewApiKeyClientWithBaseURI(endpoint string) ApiKeyClient {
	return ApiKeyClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

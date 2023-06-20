package querypacks

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QueryPacksClient struct {
	Client  autorest.Client
	baseUri string
}

func NewQueryPacksClientWithBaseURI(endpoint string) QueryPacksClient {
	return QueryPacksClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

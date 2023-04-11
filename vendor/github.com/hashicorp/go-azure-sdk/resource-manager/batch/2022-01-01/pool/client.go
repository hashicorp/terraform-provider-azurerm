package pool

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PoolClient struct {
	Client  autorest.Client
	baseUri string
}

func NewPoolClientWithBaseURI(endpoint string) PoolClient {
	return PoolClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

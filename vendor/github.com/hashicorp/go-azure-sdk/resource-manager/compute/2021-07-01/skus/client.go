package skus

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SkusClient struct {
	Client  autorest.Client
	baseUri string
}

func NewSkusClientWithBaseURI(endpoint string) SkusClient {
	return SkusClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

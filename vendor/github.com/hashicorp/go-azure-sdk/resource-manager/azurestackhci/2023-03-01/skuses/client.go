package skuses

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SkusesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewSkusesClientWithBaseURI(endpoint string) SkusesClient {
	return SkusesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

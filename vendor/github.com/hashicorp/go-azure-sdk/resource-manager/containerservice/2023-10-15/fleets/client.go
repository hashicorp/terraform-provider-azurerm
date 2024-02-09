package fleets

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FleetsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewFleetsClientWithBaseURI(endpoint string) FleetsClient {
	return FleetsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

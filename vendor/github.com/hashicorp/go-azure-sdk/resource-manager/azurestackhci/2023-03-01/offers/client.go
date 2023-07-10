package offers

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OffersClient struct {
	Client  autorest.Client
	baseUri string
}

func NewOffersClientWithBaseURI(endpoint string) OffersClient {
	return OffersClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

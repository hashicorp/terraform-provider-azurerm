package views

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ViewsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewViewsClientWithBaseURI(endpoint string) ViewsClient {
	return ViewsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

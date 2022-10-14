package links

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LinksClient struct {
	Client  autorest.Client
	baseUri string
}

func NewLinksClientWithBaseURI(endpoint string) LinksClient {
	return LinksClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

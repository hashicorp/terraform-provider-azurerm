package site

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SiteClient struct {
	Client  autorest.Client
	baseUri string
}

func NewSiteClientWithBaseURI(endpoint string) SiteClient {
	return SiteClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

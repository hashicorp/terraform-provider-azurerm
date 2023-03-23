package exports

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExportsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewExportsClientWithBaseURI(endpoint string) ExportsClient {
	return ExportsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

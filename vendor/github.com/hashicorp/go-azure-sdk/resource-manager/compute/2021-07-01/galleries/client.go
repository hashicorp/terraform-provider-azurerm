package galleries

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GalleriesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewGalleriesClientWithBaseURI(endpoint string) GalleriesClient {
	return GalleriesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

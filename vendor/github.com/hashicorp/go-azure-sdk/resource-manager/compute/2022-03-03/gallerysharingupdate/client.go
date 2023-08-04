package gallerysharingupdate

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GallerySharingUpdateClient struct {
	Client  autorest.Client
	baseUri string
}

func NewGallerySharingUpdateClientWithBaseURI(endpoint string) GallerySharingUpdateClient {
	return GallerySharingUpdateClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

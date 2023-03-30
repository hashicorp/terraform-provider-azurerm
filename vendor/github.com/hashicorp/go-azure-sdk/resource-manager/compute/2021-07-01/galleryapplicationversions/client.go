package galleryapplicationversions

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GalleryApplicationVersionsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewGalleryApplicationVersionsClientWithBaseURI(endpoint string) GalleryApplicationVersionsClient {
	return GalleryApplicationVersionsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

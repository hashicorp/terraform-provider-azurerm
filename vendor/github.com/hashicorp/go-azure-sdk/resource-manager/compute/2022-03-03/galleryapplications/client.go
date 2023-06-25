package galleryapplications

import "github.com/Azure/go-autorest/autorest"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GalleryApplicationsClient struct {
	Client  autorest.Client
	baseUri string
}

func NewGalleryApplicationsClientWithBaseURI(endpoint string) GalleryApplicationsClient {
	return GalleryApplicationsClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

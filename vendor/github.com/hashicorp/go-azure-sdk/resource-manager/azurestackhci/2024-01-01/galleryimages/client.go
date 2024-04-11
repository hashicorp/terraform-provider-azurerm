package galleryimages

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GalleryImagesClient struct {
	Client *resourcemanager.Client
}

func NewGalleryImagesClientWithBaseURI(sdkApi sdkEnv.Api) (*GalleryImagesClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "galleryimages", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating GalleryImagesClient: %+v", err)
	}

	return &GalleryImagesClient{
		Client: client,
	}, nil
}

package galleryimageversions

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GalleryImageVersionsClient struct {
	Client *resourcemanager.Client
}

func NewGalleryImageVersionsClientWithBaseURI(sdkApi sdkEnv.Api) (*GalleryImageVersionsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "galleryimageversions", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating GalleryImageVersionsClient: %+v", err)
	}

	return &GalleryImageVersionsClient{
		Client: client,
	}, nil
}

package marketplacegalleryimages

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MarketplaceGalleryImagesClient struct {
	Client *resourcemanager.Client
}

func NewMarketplaceGalleryImagesClientWithBaseURI(sdkApi sdkEnv.Api) (*MarketplaceGalleryImagesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "marketplacegalleryimages", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating MarketplaceGalleryImagesClient: %+v", err)
	}

	return &MarketplaceGalleryImagesClient{
		Client: client,
	}, nil
}

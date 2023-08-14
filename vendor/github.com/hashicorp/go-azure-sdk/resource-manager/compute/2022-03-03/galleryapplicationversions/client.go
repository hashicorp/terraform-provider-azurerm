package galleryapplicationversions

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GalleryApplicationVersionsClient struct {
	Client *resourcemanager.Client
}

func NewGalleryApplicationVersionsClientWithBaseURI(sdkApi sdkEnv.Api) (*GalleryApplicationVersionsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "galleryapplicationversions", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating GalleryApplicationVersionsClient: %+v", err)
	}

	return &GalleryApplicationVersionsClient{
		Client: client,
	}, nil
}

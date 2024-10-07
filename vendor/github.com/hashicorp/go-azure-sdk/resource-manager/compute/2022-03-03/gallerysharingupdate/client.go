package gallerysharingupdate

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GallerySharingUpdateClient struct {
	Client *resourcemanager.Client
}

func NewGallerySharingUpdateClientWithBaseURI(sdkApi sdkEnv.Api) (*GallerySharingUpdateClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "gallerysharingupdate", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating GallerySharingUpdateClient: %+v", err)
	}

	return &GallerySharingUpdateClient{
		Client: client,
	}, nil
}

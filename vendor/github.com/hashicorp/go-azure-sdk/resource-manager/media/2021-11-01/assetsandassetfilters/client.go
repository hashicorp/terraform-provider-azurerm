package assetsandassetfilters

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AssetsAndAssetFiltersClient struct {
	Client *resourcemanager.Client
}

func NewAssetsAndAssetFiltersClientWithBaseURI(sdkApi sdkEnv.Api) (*AssetsAndAssetFiltersClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "assetsandassetfilters", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating AssetsAndAssetFiltersClient: %+v", err)
	}

	return &AssetsAndAssetFiltersClient{
		Client: client,
	}, nil
}

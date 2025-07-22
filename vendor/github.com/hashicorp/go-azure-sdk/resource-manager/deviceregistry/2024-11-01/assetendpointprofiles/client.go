package assetendpointprofiles

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AssetEndpointProfilesClient struct {
	Client *resourcemanager.Client
}

func NewAssetEndpointProfilesClientWithBaseURI(sdkApi sdkEnv.Api) (*AssetEndpointProfilesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "assetendpointprofiles", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating AssetEndpointProfilesClient: %+v", err)
	}

	return &AssetEndpointProfilesClient{
		Client: client,
	}, nil
}

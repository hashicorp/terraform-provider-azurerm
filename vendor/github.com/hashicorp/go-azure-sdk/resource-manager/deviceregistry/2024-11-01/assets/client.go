package assets

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AssetsClient struct {
	Client *resourcemanager.Client
}

func NewAssetsClientWithBaseURI(sdkApi sdkEnv.Api) (*AssetsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "assets", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating AssetsClient: %+v", err)
	}

	return &AssetsClient{
		Client: client,
	}, nil
}

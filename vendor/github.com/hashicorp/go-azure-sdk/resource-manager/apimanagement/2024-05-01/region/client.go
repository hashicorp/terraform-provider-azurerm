package region

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RegionClient struct {
	Client *resourcemanager.Client
}

func NewRegionClientWithBaseURI(sdkApi sdkEnv.Api) (*RegionClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "region", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating RegionClient: %+v", err)
	}

	return &RegionClient{
		Client: client,
	}, nil
}

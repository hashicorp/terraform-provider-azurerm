package fabriccapacities

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FabricCapacitiesClient struct {
	Client *resourcemanager.Client
}

func NewFabricCapacitiesClientWithBaseURI(sdkApi sdkEnv.Api) (*FabricCapacitiesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "fabriccapacities", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating FabricCapacitiesClient: %+v", err)
	}

	return &FabricCapacitiesClient{
		Client: client,
	}, nil
}

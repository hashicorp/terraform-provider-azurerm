package fleetupdatestrategies

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FleetUpdateStrategiesClient struct {
	Client *resourcemanager.Client
}

func NewFleetUpdateStrategiesClientWithBaseURI(sdkApi sdkEnv.Api) (*FleetUpdateStrategiesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "fleetupdatestrategies", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating FleetUpdateStrategiesClient: %+v", err)
	}

	return &FleetUpdateStrategiesClient{
		Client: client,
	}, nil
}

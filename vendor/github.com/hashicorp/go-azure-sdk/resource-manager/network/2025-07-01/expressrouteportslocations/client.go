package expressrouteportslocations

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRoutePortsLocationsClient struct {
	Client *resourcemanager.Client
}

func NewExpressRoutePortsLocationsClientWithBaseURI(sdkApi sdkEnv.Api) (*ExpressRoutePortsLocationsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "expressrouteportslocations", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ExpressRoutePortsLocationsClient: %+v", err)
	}

	return &ExpressRoutePortsLocationsClient{
		Client: client,
	}, nil
}

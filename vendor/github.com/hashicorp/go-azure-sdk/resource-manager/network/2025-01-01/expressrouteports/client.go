package expressrouteports

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRoutePortsClient struct {
	Client *resourcemanager.Client
}

func NewExpressRoutePortsClientWithBaseURI(sdkApi sdkEnv.Api) (*ExpressRoutePortsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "expressrouteports", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ExpressRoutePortsClient: %+v", err)
	}

	return &ExpressRoutePortsClient{
		Client: client,
	}, nil
}

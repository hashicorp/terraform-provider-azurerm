package expressrouteconnections

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteConnectionsClient struct {
	Client *resourcemanager.Client
}

func NewExpressRouteConnectionsClientWithBaseURI(api environments.Api) (*ExpressRouteConnectionsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "expressrouteconnections", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ExpressRouteConnectionsClient: %+v", err)
	}

	return &ExpressRouteConnectionsClient{
		Client: client,
	}, nil
}

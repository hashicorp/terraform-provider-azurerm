package routetables

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RouteTablesClient struct {
	Client *resourcemanager.Client
}

func NewRouteTablesClientWithBaseURI(sdkApi sdkEnv.Api) (*RouteTablesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "routetables", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating RouteTablesClient: %+v", err)
	}

	return &RouteTablesClient{
		Client: client,
	}, nil
}

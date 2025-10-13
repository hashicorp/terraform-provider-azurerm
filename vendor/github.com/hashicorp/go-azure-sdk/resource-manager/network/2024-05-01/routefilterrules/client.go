package routefilterrules

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RouteFilterRulesClient struct {
	Client *resourcemanager.Client
}

func NewRouteFilterRulesClientWithBaseURI(sdkApi sdkEnv.Api) (*RouteFilterRulesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "routefilterrules", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating RouteFilterRulesClient: %+v", err)
	}

	return &RouteFilterRulesClient{
		Client: client,
	}, nil
}

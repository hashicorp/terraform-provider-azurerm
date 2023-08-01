package expressroutecrossconnectionroutetable

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteCrossConnectionRouteTableClient struct {
	Client *resourcemanager.Client
}

func NewExpressRouteCrossConnectionRouteTableClientWithBaseURI(sdkApi sdkEnv.Api) (*ExpressRouteCrossConnectionRouteTableClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "expressroutecrossconnectionroutetable", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ExpressRouteCrossConnectionRouteTableClient: %+v", err)
	}

	return &ExpressRouteCrossConnectionRouteTableClient{
		Client: client,
	}, nil
}

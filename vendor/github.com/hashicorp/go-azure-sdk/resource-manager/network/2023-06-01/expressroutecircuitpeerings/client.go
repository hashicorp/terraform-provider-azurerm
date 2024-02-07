package expressroutecircuitpeerings

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteCircuitPeeringsClient struct {
	Client *resourcemanager.Client
}

func NewExpressRouteCircuitPeeringsClientWithBaseURI(sdkApi sdkEnv.Api) (*ExpressRouteCircuitPeeringsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "expressroutecircuitpeerings", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ExpressRouteCircuitPeeringsClient: %+v", err)
	}

	return &ExpressRouteCircuitPeeringsClient{
		Client: client,
	}, nil
}

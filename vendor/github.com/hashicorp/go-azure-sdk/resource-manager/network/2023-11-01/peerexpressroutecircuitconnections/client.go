package peerexpressroutecircuitconnections

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PeerExpressRouteCircuitConnectionsClient struct {
	Client *resourcemanager.Client
}

func NewPeerExpressRouteCircuitConnectionsClientWithBaseURI(sdkApi sdkEnv.Api) (*PeerExpressRouteCircuitConnectionsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "peerexpressroutecircuitconnections", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating PeerExpressRouteCircuitConnectionsClient: %+v", err)
	}

	return &PeerExpressRouteCircuitConnectionsClient{
		Client: client,
	}, nil
}

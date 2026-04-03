package virtualnetworkpeerings

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualNetworkPeeringsClient struct {
	Client *resourcemanager.Client
}

func NewVirtualNetworkPeeringsClientWithBaseURI(sdkApi sdkEnv.Api) (*VirtualNetworkPeeringsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "virtualnetworkpeerings", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating VirtualNetworkPeeringsClient: %+v", err)
	}

	return &VirtualNetworkPeeringsClient{
		Client: client,
	}, nil
}

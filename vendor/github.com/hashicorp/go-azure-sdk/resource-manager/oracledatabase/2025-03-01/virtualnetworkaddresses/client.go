package virtualnetworkaddresses

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualNetworkAddressesClient struct {
	Client *resourcemanager.Client
}

func NewVirtualNetworkAddressesClientWithBaseURI(sdkApi sdkEnv.Api) (*VirtualNetworkAddressesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "virtualnetworkaddresses", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating VirtualNetworkAddressesClient: %+v", err)
	}

	return &VirtualNetworkAddressesClient{
		Client: client,
	}, nil
}

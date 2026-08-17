package virtualnetworktaps

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualNetworkTapsClient struct {
	Client *resourcemanager.Client
}

func NewVirtualNetworkTapsClientWithBaseURI(sdkApi sdkEnv.Api) (*VirtualNetworkTapsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "virtualnetworktaps", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating VirtualNetworkTapsClient: %+v", err)
	}

	return &VirtualNetworkTapsClient{
		Client: client,
	}, nil
}

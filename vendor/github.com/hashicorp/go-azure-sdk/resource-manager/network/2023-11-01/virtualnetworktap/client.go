package virtualnetworktap

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualNetworkTapClient struct {
	Client *resourcemanager.Client
}

func NewVirtualNetworkTapClientWithBaseURI(sdkApi sdkEnv.Api) (*VirtualNetworkTapClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "virtualnetworktap", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating VirtualNetworkTapClient: %+v", err)
	}

	return &VirtualNetworkTapClient{
		Client: client,
	}, nil
}

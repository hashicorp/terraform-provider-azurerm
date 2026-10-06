package networkgateways

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkGatewaysClient struct {
	Client *resourcemanager.Client
}

func NewNetworkGatewaysClientWithBaseURI(sdkApi sdkEnv.Api) (*NetworkGatewaysClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "networkgateways", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating NetworkGatewaysClient: %+v", err)
	}

	return &NetworkGatewaysClient{
		Client: client,
	}, nil
}

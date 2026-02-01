package localnetworkgateways

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LocalNetworkGatewaysClient struct {
	Client *resourcemanager.Client
}

func NewLocalNetworkGatewaysClientWithBaseURI(sdkApi sdkEnv.Api) (*LocalNetworkGatewaysClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "localnetworkgateways", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating LocalNetworkGatewaysClient: %+v", err)
	}

	return &LocalNetworkGatewaysClient{
		Client: client,
	}, nil
}

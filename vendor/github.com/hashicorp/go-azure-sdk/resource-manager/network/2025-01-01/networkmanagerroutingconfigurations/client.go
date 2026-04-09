package networkmanagerroutingconfigurations

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkManagerRoutingConfigurationsClient struct {
	Client *resourcemanager.Client
}

func NewNetworkManagerRoutingConfigurationsClientWithBaseURI(sdkApi sdkEnv.Api) (*NetworkManagerRoutingConfigurationsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "networkmanagerroutingconfigurations", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating NetworkManagerRoutingConfigurationsClient: %+v", err)
	}

	return &NetworkManagerRoutingConfigurationsClient{
		Client: client,
	}, nil
}

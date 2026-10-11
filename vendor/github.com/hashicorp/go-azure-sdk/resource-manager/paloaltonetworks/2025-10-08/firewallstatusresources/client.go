package firewallstatusresources

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FirewallStatusResourcesClient struct {
	Client *resourcemanager.Client
}

func NewFirewallStatusResourcesClientWithBaseURI(sdkApi sdkEnv.Api) (*FirewallStatusResourcesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "firewallstatusresources", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating FirewallStatusResourcesClient: %+v", err)
	}

	return &FirewallStatusResourcesClient{
		Client: client,
	}, nil
}

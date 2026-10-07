package firewallresources

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FirewallResourcesClient struct {
	Client *resourcemanager.Client
}

func NewFirewallResourcesClientWithBaseURI(sdkApi sdkEnv.Api) (*FirewallResourcesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "firewallresources", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating FirewallResourcesClient: %+v", err)
	}

	return &FirewallResourcesClient{
		Client: client,
	}, nil
}

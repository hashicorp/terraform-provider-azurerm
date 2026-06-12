package firewallrules

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FirewallRulesClient struct {
	Client *resourcemanager.Client
}

func NewFirewallRulesClientWithBaseURI(sdkApi sdkEnv.Api) (*FirewallRulesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "firewallrules", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating FirewallRulesClient: %+v", err)
	}

	return &FirewallRulesClient{
		Client: client,
	}, nil
}

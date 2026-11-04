package outboundfirewallrules

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OutboundFirewallRulesClient struct {
	Client *resourcemanager.Client
}

func NewOutboundFirewallRulesClientWithBaseURI(sdkApi sdkEnv.Api) (*OutboundFirewallRulesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "outboundfirewallrules", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating OutboundFirewallRulesClient: %+v", err)
	}

	return &OutboundFirewallRulesClient{
		Client: client,
	}, nil
}

package ipfirewallrules

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IPFirewallRulesClient struct {
	Client *resourcemanager.Client
}

func NewIPFirewallRulesClientWithBaseURI(sdkApi sdkEnv.Api) (*IPFirewallRulesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "ipfirewallrules", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating IPFirewallRulesClient: %+v", err)
	}

	return &IPFirewallRulesClient{
		Client: client,
	}, nil
}

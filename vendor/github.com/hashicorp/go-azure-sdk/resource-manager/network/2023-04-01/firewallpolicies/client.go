package firewallpolicies

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FirewallPoliciesClient struct {
	Client *resourcemanager.Client
}

func NewFirewallPoliciesClientWithBaseURI(api environments.Api) (*FirewallPoliciesClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "firewallpolicies", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating FirewallPoliciesClient: %+v", err)
	}

	return &FirewallPoliciesClient{
		Client: client,
	}, nil
}

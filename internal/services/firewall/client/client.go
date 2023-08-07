// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	network_2023_04_01 "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
)

type Client struct {
	*network_2023_04_01.Client

	FirewallPolicyClient          *network.FirewallPoliciesClient
	FirewallPolicyRuleGroupClient *network.FirewallPolicyRuleCollectionGroupsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	client, err := network_2023_04_01.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building clients for Network: %+v", err)
	}

	policyClient := network.NewFirewallPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&policyClient.Client, o.ResourceManagerAuthorizer)

	policyRuleGroupClient := network.NewFirewallPolicyRuleCollectionGroupsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&policyRuleGroupClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		Client:                        client,
		FirewallPolicyClient:          &policyClient,
		FirewallPolicyRuleGroupClient: &policyRuleGroupClient,
	}, nil
}

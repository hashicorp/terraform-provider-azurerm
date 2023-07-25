// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
)

type Client struct {
	AzureFirewallsClient          *network.AzureFirewallsClient
	FirewallPolicyClient          *network.FirewallPoliciesClient
	FirewallPolicyRuleGroupClient *network.FirewallPolicyRuleCollectionGroupsClient
}

func NewClient(o *common.ClientOptions) *Client {
	firewallsClient := network.NewAzureFirewallsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&firewallsClient.Client, o.ResourceManagerAuthorizer)

	policyClient := network.NewFirewallPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&policyClient.Client, o.ResourceManagerAuthorizer)

	policyRuleGroupClient := network.NewFirewallPolicyRuleCollectionGroupsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&policyRuleGroupClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AzureFirewallsClient:          &firewallsClient,
		FirewallPolicyClient:          &policyClient,
		FirewallPolicyRuleGroupClient: &policyRuleGroupClient,
	}
}

package client

import (
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2021-08-01/network"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	network2 "github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall/sdk/2022-01-01/network"
)

type Client struct {
	AzureFirewallsClient          *network.AzureFirewallsClient
	FirewallPolicyClient          *network2.FirewallPoliciesClient
	FirewallPolicyRuleGroupClient *network.FirewallPolicyRuleCollectionGroupsClient
}

func NewClient(o *common.ClientOptions) *Client {
	firewallsClient := network.NewAzureFirewallsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&firewallsClient.Client, o.ResourceManagerAuthorizer)

	policyClient := network2.NewFirewallPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&policyClient.Client, o.ResourceManagerAuthorizer)

	policyRuleGroupClient := network.NewFirewallPolicyRuleCollectionGroupsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&policyRuleGroupClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AzureFirewallsClient:          &firewallsClient,
		FirewallPolicyClient:          &policyClient,
		FirewallPolicyRuleGroupClient: &policyRuleGroupClient,
	}
}

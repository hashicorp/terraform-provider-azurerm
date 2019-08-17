package webapplicationfirewallpolicies

import (
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-04-01/network"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	WebApplicationFirewallPoliciesClient *network.WebApplicationFirewallPoliciesClient
}

func BuildClient(o *common.ClientOptions) *Client {

	WebApplicationFirewallPoliciesClient := network.NewWebApplicationFirewallPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&WebApplicationFirewallPoliciesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		WebApplicationFirewallPoliciesClient: &WebApplicationFirewallPoliciesClient,
	}
}

package client

import (
	networkLegacy "github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-03-01/network"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	LoadBalancersClient      *networkLegacy.LoadBalancersClient
	LoadBalancingRulesClient *networkLegacy.LoadBalancerLoadBalancingRulesClient
}

func NewClient(o *common.ClientOptions) *Client {
	loadBalancersClient := networkLegacy.NewLoadBalancersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&loadBalancersClient.Client, o.ResourceManagerAuthorizer)

	loadBalancingRulesClient := networkLegacy.NewLoadBalancerLoadBalancingRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&loadBalancingRulesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		LoadBalancersClient:      &loadBalancersClient,
		LoadBalancingRulesClient: &loadBalancingRulesClient,
	}
}

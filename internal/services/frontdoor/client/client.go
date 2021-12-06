package client

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/frontdoor/sdk/2020-04-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/frontdoor/sdk/2020-05-01/frontdoors"
)

type Client struct {
	FrontDoorsClient       *frontdoors.FrontDoorsClient
	FrontDoorsPolicyClient *webapplicationfirewallpolicies.WebApplicationFirewallPoliciesClient
}

func NewClient(o *common.ClientOptions) *Client {
	frontDoorsClient := frontdoors.NewFrontDoorsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&frontDoorsClient.Client, o.ResourceManagerAuthorizer)

	frontDoorsPolicyClient := webapplicationfirewallpolicies.NewWebApplicationFirewallPoliciesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&frontDoorsPolicyClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		FrontDoorsClient:       &frontDoorsClient,
		FrontDoorsPolicyClient: &frontDoorsPolicyClient,
	}
}

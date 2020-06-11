package client

import (
	"github.com/Azure/azure-sdk-for-go/services/frontdoor/mgmt/2020-01-01/frontdoor"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	FrontDoorsClient         *frontdoor.FrontDoorsClient
	FrontDoorsFrontendClient *frontdoor.FrontendEndpointsClient
	FrontDoorsPolicyClient   *frontdoor.PoliciesClient
}

func NewClient(o *common.ClientOptions) *Client {
	frontDoorsClient := frontdoor.NewFrontDoorsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&frontDoorsClient.Client, o.ResourceManagerAuthorizer)

	frontDoorsFrontendClient := frontdoor.NewFrontendEndpointsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&frontDoorsFrontendClient.Client, o.ResourceManagerAuthorizer)

	frontDoorsPolicyClient := frontdoor.NewPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&frontDoorsPolicyClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		FrontDoorsClient:         &frontDoorsClient,
		FrontDoorsFrontendClient: &frontDoorsFrontendClient,
		FrontDoorsPolicyClient:   &frontDoorsPolicyClient,
	}
}

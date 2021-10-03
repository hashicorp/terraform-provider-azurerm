package client

import (
	"github.com/Azure/azure-sdk-for-go/services/frontdoor/mgmt/2020-05-01/frontdoor"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	FrontDoorsClient             *frontdoor.FrontDoorsClient
	FrontDoorsFrontendClient     *frontdoor.FrontendEndpointsClient
	FrontDoorsPolicyClient       *frontdoor.PoliciesClient
	FrontDoorsRulesEnginesClient *frontdoor.RulesEnginesClient
}

func NewClient(o *common.ClientOptions) *Client {
	frontDoorsClient := frontdoor.NewFrontDoorsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&frontDoorsClient.Client, o.ResourceManagerAuthorizer)

	frontDoorsFrontendClient := frontdoor.NewFrontendEndpointsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&frontDoorsFrontendClient.Client, o.ResourceManagerAuthorizer)

	frontDoorsRulesEnginesClient := frontdoor.NewRulesEnginesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&frontDoorsRulesEnginesClient.Client, o.ResourceManagerAuthorizer)

	frontDoorsPolicyClient := frontdoor.NewPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&frontDoorsPolicyClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		FrontDoorsClient:             &frontDoorsClient,
		FrontDoorsFrontendClient:     &frontDoorsFrontendClient,
		FrontDoorsPolicyClient:       &frontDoorsPolicyClient,
		FrontDoorsRulesEnginesClient: &frontDoorsRulesEnginesClient,
	}
}

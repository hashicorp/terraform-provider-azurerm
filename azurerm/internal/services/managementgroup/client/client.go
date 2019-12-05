package client

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/resources/mgmt/2018-03-01-preview/managementgroups"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	GroupsClient       *managementgroups.Client
	SubscriptionClient *managementgroups.SubscriptionsClient
}

func NewClient(o *common.ClientOptions) *Client {
	GroupsClient := managementgroups.NewClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&GroupsClient.Client, o.ResourceManagerAuthorizer)

	SubscriptionClient := managementgroups.NewSubscriptionsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&SubscriptionClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		GroupsClient:       &GroupsClient,
		SubscriptionClient: &SubscriptionClient,
	}
}

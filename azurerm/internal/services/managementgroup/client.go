package managementgroup

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/resources/mgmt/2018-03-01-preview/managementgroups"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	GroupsClient       managementgroups.Client
	SubscriptionClient managementgroups.SubscriptionsClient
}

func BuildClient(o *common.ClientOptions) *Client {
	c := Client{}

	c.GroupsClient = managementgroups.NewClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&c.GroupsClient.Client, o.ResourceManagerAuthorizer)

	c.SubscriptionClient = managementgroups.NewSubscriptionsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&c.SubscriptionClient.Client, o.ResourceManagerAuthorizer)

	return &c
}

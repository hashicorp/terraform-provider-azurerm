package servicebus

import (
	"github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2017-04-01/servicebus"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	QueuesClient            servicebus.QueuesClient
	NamespacesClient        servicebus.NamespacesClient
	TopicsClient            servicebus.TopicsClient
	SubscriptionsClient     servicebus.SubscriptionsClient
	SubscriptionRulesClient servicebus.RulesClient
}

func BuildClient(o *common.ClientOptions) *Client {
	c := Client{}

	c.QueuesClient = servicebus.NewQueuesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.QueuesClient.Client, o.ResourceManagerAuthorizer)

	c.NamespacesClient = servicebus.NewNamespacesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.NamespacesClient.Client, o.ResourceManagerAuthorizer)

	c.TopicsClient = servicebus.NewTopicsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.TopicsClient.Client, o.ResourceManagerAuthorizer)

	c.SubscriptionsClient = servicebus.NewSubscriptionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.SubscriptionsClient.Client, o.ResourceManagerAuthorizer)

	c.SubscriptionRulesClient = servicebus.NewRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.SubscriptionRulesClient.Client, o.ResourceManagerAuthorizer)

	return &c
}

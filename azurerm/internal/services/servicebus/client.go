package servicebus

import (
	"github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2017-04-01/servicebus"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	QueuesClient            *servicebus.QueuesClient
	NamespacesClient        *servicebus.NamespacesClient
	TopicsClient            *servicebus.TopicsClient
	SubscriptionsClient     *servicebus.SubscriptionsClient
	SubscriptionRulesClient *servicebus.RulesClient
}

func BuildClient(o *common.ClientOptions) *Client {

	QueuesClient := servicebus.NewQueuesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&QueuesClient.Client, o.ResourceManagerAuthorizer)

	NamespacesClient := servicebus.NewNamespacesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&NamespacesClient.Client, o.ResourceManagerAuthorizer)

	TopicsClient := servicebus.NewTopicsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&TopicsClient.Client, o.ResourceManagerAuthorizer)

	SubscriptionsClient := servicebus.NewSubscriptionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&SubscriptionsClient.Client, o.ResourceManagerAuthorizer)

	SubscriptionRulesClient := servicebus.NewRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&SubscriptionRulesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		QueuesClient:            &QueuesClient,
		NamespacesClient:        &NamespacesClient,
		TopicsClient:            &TopicsClient,
		SubscriptionsClient:     &SubscriptionsClient,
		SubscriptionRulesClient: &SubscriptionRulesClient,
	}
}

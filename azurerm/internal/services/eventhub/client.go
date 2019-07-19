package eventhub

import (
	"github.com/Azure/azure-sdk-for-go/services/eventhub/mgmt/2017-04-01/eventhub"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ConsumerGroupClient eventhub.ConsumerGroupsClient
	EventHubsClient     eventhub.EventHubsClient
	NamespacesClient    eventhub.NamespacesClient
}

func BuildClient(o *common.ClientOptions) *Client {
	c := Client{}

	c.EventHubsClient = eventhub.NewEventHubsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.EventHubsClient.Client, o.ResourceManagerAuthorizer)

	c.ConsumerGroupClient = eventhub.NewConsumerGroupsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ConsumerGroupClient.Client, o.ResourceManagerAuthorizer)

	c.NamespacesClient = eventhub.NewNamespacesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.NamespacesClient.Client, o.ResourceManagerAuthorizer)

	return &c
}

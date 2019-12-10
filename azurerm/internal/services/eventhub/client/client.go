package client

import (
	"github.com/Azure/azure-sdk-for-go/services/eventhub/mgmt/2017-04-01/eventhub"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ConsumerGroupClient           *eventhub.ConsumerGroupsClient
	DisasterRecoveryConfigsClient *eventhub.DisasterRecoveryConfigsClient
	EventHubsClient               *eventhub.EventHubsClient
	NamespacesClient              *eventhub.NamespacesClient
}

func NewClient(o *common.ClientOptions) *Client {
	EventHubsClient := eventhub.NewEventHubsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&EventHubsClient.Client, o.ResourceManagerAuthorizer)

	DisasterRecoveryConfigsClient := eventhub.NewDisasterRecoveryConfigsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&DisasterRecoveryConfigsClient.Client, o.ResourceManagerAuthorizer)

	ConsumerGroupClient := eventhub.NewConsumerGroupsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ConsumerGroupClient.Client, o.ResourceManagerAuthorizer)

	NamespacesClient := eventhub.NewNamespacesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&NamespacesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ConsumerGroupClient:           &ConsumerGroupClient,
		DisasterRecoveryConfigsClient: &DisasterRecoveryConfigsClient,
		EventHubsClient:               &EventHubsClient,
		NamespacesClient:              &NamespacesClient,
	}
}

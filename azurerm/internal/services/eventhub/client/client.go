package client

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/eventhub/mgmt/2018-01-01-preview/eventhub"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ClusterClient                 *eventhub.ClustersClient
	ConsumerGroupClient           *eventhub.ConsumerGroupsClient
	DisasterRecoveryConfigsClient *eventhub.DisasterRecoveryConfigsClient
	EventHubsClient               *eventhub.EventHubsClient
	NamespacesClient              *eventhub.NamespacesClient
}

func NewClient(o *common.ClientOptions) *Client {
	ClustersClient := eventhub.NewClustersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ClustersClient.Client, o.ResourceManagerAuthorizer)

	EventHubsClient := eventhub.NewEventHubsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&EventHubsClient.Client, o.ResourceManagerAuthorizer)

	DisasterRecoveryConfigsClient := eventhub.NewDisasterRecoveryConfigsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&DisasterRecoveryConfigsClient.Client, o.ResourceManagerAuthorizer)

	ConsumerGroupClient := eventhub.NewConsumerGroupsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ConsumerGroupClient.Client, o.ResourceManagerAuthorizer)

	NamespacesClient := eventhub.NewNamespacesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&NamespacesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ClusterClient:                 &ClustersClient,
		ConsumerGroupClient:           &ConsumerGroupClient,
		DisasterRecoveryConfigsClient: &DisasterRecoveryConfigsClient,
		EventHubsClient:               &EventHubsClient,
		NamespacesClient:              &NamespacesClient,
	}
}

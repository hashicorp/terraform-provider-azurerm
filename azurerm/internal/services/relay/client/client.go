package client

import (
	"github.com/Azure/azure-sdk-for-go/services/relay/mgmt/2017-04-01/relay"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	HybridConnectionsClient *relay.HybridConnectionsClient
	NamespacesClient        *relay.NamespacesClient
}

func NewClient(o *common.ClientOptions) *Client {
	hybridConnectionsClient := relay.NewHybridConnectionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&hybridConnectionsClient.Client, o.ResourceManagerAuthorizer)

	namespacesClient := relay.NewNamespacesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&namespacesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		HybridConnectionsClient: &hybridConnectionsClient,
		NamespacesClient:        &namespacesClient,
	}
}

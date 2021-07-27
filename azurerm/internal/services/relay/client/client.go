package client

import (
	"github.com/Azure/azure-sdk-for-go/services/relay/mgmt/2017-04-01/relay"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/relay/sdk/namespaces"
)

type Client struct {
	HybridConnectionsClient *relay.HybridConnectionsClient
	NamespacesClient        *namespaces.NamespacesClient
}

func NewClient(o *common.ClientOptions) *Client {
	hybridConnectionsClient := relay.NewHybridConnectionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&hybridConnectionsClient.Client, o.ResourceManagerAuthorizer)

	namespacesClient := namespaces.NewNamespacesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&namespacesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		HybridConnectionsClient: &hybridConnectionsClient,
		NamespacesClient:        &namespacesClient,
	}
}

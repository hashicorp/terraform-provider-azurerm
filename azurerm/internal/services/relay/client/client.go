package client

import (
	"github.com/Azure/azure-sdk-for-go/services/relay/mgmt/2017-04-01/relay"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/relay/sdk/hybridconnections"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/relay/sdk/namespaces"
)

type Client struct {
	HybridConnectionsClient *hybridconnections.HybridConnectionsClient
	NamespacesClient        *namespaces.NamespacesClient

	// TODO: work through and enable removing this
	Track1NamespacesClient *relay.NamespacesClient
}

func NewClient(o *common.ClientOptions) *Client {
	hybridConnectionsClient := hybridconnections.NewHybridConnectionsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&hybridConnectionsClient.Client, o.ResourceManagerAuthorizer)

	namespacesClient := namespaces.NewNamespacesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&namespacesClient.Client, o.ResourceManagerAuthorizer)

	track1NamespacesClient := relay.NewNamespacesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&track1NamespacesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		HybridConnectionsClient: &hybridConnectionsClient,
		NamespacesClient:        &namespacesClient,

		Track1NamespacesClient: &track1NamespacesClient,
	}
}

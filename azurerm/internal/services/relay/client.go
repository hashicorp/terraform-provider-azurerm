package relay

import (
	"github.com/Azure/azure-sdk-for-go/services/relay/mgmt/2017-04-01/relay"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	NamespacesClient        *relay.NamespacesClient
	HybridConnectionsClient *relay.HybridConnectionsClient
}

func BuildClient(o *common.ClientOptions) *Client {
	NamespacesClient := relay.NewNamespacesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&NamespacesClient.Client, o.ResourceManagerAuthorizer)

	HybridConnectionsClient := relay.NewHybridConnectionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&HybridConnectionsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		NamespacesClient:        &NamespacesClient,
		HybridConnectionsClient: &HybridConnectionsClient,
	}
}

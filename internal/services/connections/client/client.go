package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2016-06-01/connections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2016-06-01/managedapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ConnectionsClient *connections.ConnectionsClient
	ManagedApisClient *managedapis.ManagedAPIsClient
}

func NewClient(o *common.ClientOptions) *Client {
	connectionsClient := connections.NewConnectionsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&connectionsClient.Client, o.ResourceManagerAuthorizer)

	managedApisClient := managedapis.NewManagedAPIsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&managedApisClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ConnectionsClient: &connectionsClient,
		ManagedApisClient: &managedApisClient,
	}
}

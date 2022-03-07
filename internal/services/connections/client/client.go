package client

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/connections/sdk/2016-06-01/connections"
)

type Client struct {
	ConnectionsClient *connections.ConnectionsClient
}

func NewClient(o *common.ClientOptions) *Client {
	connectionsClient := connections.NewConnectionsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&connectionsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ConnectionsClient: &connectionsClient,
	}
}

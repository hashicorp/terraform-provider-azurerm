package client

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/analysisservices/sdk/servers"
)

type Client struct {
	ServerClient *servers.ServersClient
}

func NewClient(o *common.ClientOptions) *Client {
	serverClient := servers.NewServersClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&serverClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ServerClient: &serverClient,
	}
}

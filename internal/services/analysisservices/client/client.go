package client

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/analysisservices/sdk/2017-08-01/servers"
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

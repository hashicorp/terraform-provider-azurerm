package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/analysisservices/2017-08-01/servers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
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

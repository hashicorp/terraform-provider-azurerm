package client

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/postgresqlhsc/sdk/2020-10-05-privatepreview/servergroups"
)

type Client struct {
	ServerGroupsClient *servergroups.ServerGroupsClient
}

func NewClient(o *common.ClientOptions) *Client {
	serverGroupsClient := servergroups.NewServerGroupsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&serverGroupsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ServerGroupsClient: &serverGroupsClient,
	}
}

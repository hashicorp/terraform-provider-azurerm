package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresqlhsc/2022-11-08/roles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	RolesClient *roles.RolesClient
}

func NewClient(o *common.ClientOptions) *Client {
	rolesClient := roles.NewRolesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&rolesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		RolesClient: &rolesClient,
	}
}

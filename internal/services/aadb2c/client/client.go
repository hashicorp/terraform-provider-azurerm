package client

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/aadb2c/sdk/2021-04-01-preview/tenants"
)

type Client struct {
	TenantsClient *tenants.TenantsClient
}

func NewClient(o *common.ClientOptions) *Client {
	tenantsClient := tenants.NewTenantsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&tenantsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		TenantsClient: &tenantsClient,
	}
}

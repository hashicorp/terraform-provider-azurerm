package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/aadb2c/2021-04-01-preview/tenants"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
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

package v2021_04_01_preview

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-sdk/resource-manager/aadb2c/2021-04-01-preview/tenants"
)

type Client struct {
	Tenants *tenants.TenantsClient
}

func NewClientWithBaseURI(endpoint string, configureAuthFunc func(c *autorest.Client)) Client {

	tenantsClient := tenants.NewTenantsClientWithBaseURI(endpoint)
	configureAuthFunc(&tenantsClient.Client)

	return Client{
		Tenants: &tenantsClient,
	}
}

package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresqlhsc/2022-11-08/configurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ConfigurationsClient *configurations.ConfigurationsClient
}

func NewClient(o *common.ClientOptions) *Client {
	configurationsClient := configurations.NewConfigurationsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&configurationsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ConfigurationsClient: &configurationsClient,
	}
}

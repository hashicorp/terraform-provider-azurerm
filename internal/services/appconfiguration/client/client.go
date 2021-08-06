package client

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appconfiguration/sdk/configurationstores"
)

type Client struct {
	ConfigurationStoresClient *configurationstores.ConfigurationStoresClient
}

func NewClient(o *common.ClientOptions) *Client {
	configurationStores := configurationstores.NewConfigurationStoresClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&configurationStores.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ConfigurationStoresClient: &configurationStores,
	}
}

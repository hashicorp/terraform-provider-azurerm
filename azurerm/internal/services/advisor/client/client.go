package client

import (
	"github.com/Azure/azure-sdk-for-go/services/advisor/mgmt/2017-04-19/advisor"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ConfigurationsClient *advisor.ConfigurationsClient
}

func NewClient(o *common.ClientOptions) *Client {
	ConfigurationsClient := advisor.NewConfigurationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ConfigurationsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ConfigurationsClient: &ConfigurationsClient,
	}
}

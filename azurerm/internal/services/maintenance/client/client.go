package client

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/maintenance/mgmt/2018-06-01-preview/maintenance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ConfigurationsClient *maintenance.ConfigurationsClient
}

func NewClient(o *common.ClientOptions) *Client {
	configurationsClient := maintenance.NewConfigurationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&configurationsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ConfigurationsClient: &configurationsClient,
	}
}

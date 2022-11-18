package client

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/tombuildsstuff/kermit/sdk/automanage/2022-05-04/automanage"
)

type Client struct {
	ConfigurationProfileClient *automanage.ConfigurationProfilesClient
}

func NewClient(o *common.ClientOptions) *Client {
	configurationProfileClient := automanage.NewConfigurationProfilesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&configurationProfileClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ConfigurationProfileClient: &configurationProfileClient,
	}
}

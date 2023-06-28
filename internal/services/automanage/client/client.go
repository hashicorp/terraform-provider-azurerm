package client

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/tombuildsstuff/kermit/sdk/automanage/2022-05-04/automanage"
)

type Client struct {
	ConfigurationClient        *automanage.ConfigurationProfilesClient
	ConfigurationVersionClient *automanage.ConfigurationProfilesVersionsClient
}

func NewClient(o *common.ClientOptions) *Client {
	configurationProfileClient := automanage.NewConfigurationProfilesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	configurationVersionsClient := automanage.NewConfigurationProfilesVersionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&configurationProfileClient.Client, o.ResourceManagerAuthorizer)
	o.ConfigureClient(&configurationVersionsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ConfigurationClient:        &configurationProfileClient,
		ConfigurationVersionClient: &configurationVersionsClient,
	}
}

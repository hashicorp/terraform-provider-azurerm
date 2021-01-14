package client

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/customproviders/mgmt/2018-09-01-preview/customproviders"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	CustomProviderClient *customproviders.CustomResourceProviderClient
}

func NewClient(o *common.ClientOptions) *Client {
	CustomProviderClient := customproviders.NewCustomResourceProviderClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&CustomProviderClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		CustomProviderClient: &CustomProviderClient,
	}
}

package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/customproviders/2018-09-01-preview/customresourceprovider"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	CustomProviderClient *customresourceprovider.CustomResourceProviderClient
}

func NewClient(o *common.ClientOptions) *Client {
	CustomProviderClient := customresourceprovider.NewCustomResourceProviderClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&CustomProviderClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		CustomProviderClient: &CustomProviderClient,
	}
}

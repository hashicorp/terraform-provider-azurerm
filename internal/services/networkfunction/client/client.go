package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/networkfunction/2022-11-01/azuretrafficcollectors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AzureTrafficCollectorsClient *azuretrafficcollectors.AzureTrafficCollectorsClient
}

func NewClient(o *common.ClientOptions) *Client {
	azureTrafficCollectorsClient := azuretrafficcollectors.NewAzureTrafficCollectorsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&azureTrafficCollectorsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AzureTrafficCollectorsClient: &azureTrafficCollectorsClient,
	}
}

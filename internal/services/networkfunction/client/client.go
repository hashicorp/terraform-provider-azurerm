package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/networkfunction/2022-11-01/azuretrafficcollectors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AzureTrafficCollectorsClient *azuretrafficcollectors.AzureTrafficCollectorsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	azureTrafficCollectorsClient, err := azuretrafficcollectors.NewAzureTrafficCollectorsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, err
	}

	o.Configure(azureTrafficCollectorsClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		AzureTrafficCollectorsClient: azureTrafficCollectorsClient,
	}, nil
}

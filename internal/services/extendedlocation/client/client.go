package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/extendedlocation/2021-08-15/customlocations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	CustomLocationsClient *customlocations.CustomLocationsClient
}

func NewClient(o *common.ClientOptions) *Client {
	customLocationsClient := customlocations.NewCustomLocationsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&customLocationsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		CustomLocationsClient: &customLocationsClient,
	}
}

package client

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/powerbi/sdk/capacities"
)

type Client struct {
	CapacityClient *capacities.CapacitiesClient
}

func NewClient(o *common.ClientOptions) *Client {
	capacityClient := capacities.NewCapacitiesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&capacityClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		CapacityClient: &capacityClient,
	}
}

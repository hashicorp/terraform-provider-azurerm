package client

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/powerbi/sdk/2021-01-01/capacities"
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

package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/powerbidedicated/2021-01-01/capacities"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
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

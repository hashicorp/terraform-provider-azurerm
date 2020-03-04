package client

import (
	"github.com/Azure/azure-sdk-for-go/services/powerbidedicated/mgmt/2017-10-01/powerbidedicated"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	CapacityClient *powerbidedicated.CapacitiesClient
}

func NewClient(o *common.ClientOptions) *Client {
	capacityClient := powerbidedicated.NewCapacitiesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&capacityClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		CapacityClient: &capacityClient,
	}
}

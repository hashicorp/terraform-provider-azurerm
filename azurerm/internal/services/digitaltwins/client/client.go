package client

import (
	"github.com/Azure/azure-sdk-for-go/services/digitaltwins/mgmt/2020-10-31/digitaltwins"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	DigitalTwinClient *digitaltwins.Client
}

func NewClient(o *common.ClientOptions) *Client {
	digitalTwinClient := digitaltwins.NewClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&digitalTwinClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		DigitalTwinClient: &digitalTwinClient,
	}
}

package client

import (
	"github.com/Azure/azure-sdk-for-go/services/databoxedge/mgmt/2019-08-01/databoxedge"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	OrderClient *databoxedge.OrdersClient
}

func NewClient(o *common.ClientOptions) *Client {
	orderClient := databoxedge.NewOrdersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&orderClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		OrderClient: &orderClient,
	}
}

package client

import (
	"github.com/Azure/azure-sdk-for-go/services/databoxedge/mgmt/2020-12-01/databoxedge"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	DeviceClient *databoxedge.DevicesClient
	OrderClient  *databoxedge.OrdersClient
}

func NewClient(o *common.ClientOptions) *Client {
	deviceClient := databoxedge.NewDevicesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&deviceClient.Client, o.ResourceManagerAuthorizer)

	orderClient := databoxedge.NewOrdersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&orderClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		DeviceClient: &deviceClient,
		OrderClient:  &orderClient,
	}
}

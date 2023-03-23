package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/databoxedge/2020-12-01/devices"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databoxedge/2020-12-01/orders"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	DeviceClient *devices.DevicesClient
	OrdersClient *orders.OrdersClient
}

func NewClient(o *common.ClientOptions) *Client {
	deviceClient := devices.NewDevicesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&deviceClient.Client, o.ResourceManagerAuthorizer)

	ordersClient := orders.NewOrdersClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&ordersClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		DeviceClient: &deviceClient,
		OrdersClient: &ordersClient,
	}
}

package client

import (
	"github.com/Azure/azure-sdk-for-go/services/databoxedge/mgmt/2020-12-01/databoxedge" // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/databoxedge/2020-12-01/devices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	DeviceClient *devices.DevicesClient
	OrderClient  *databoxedge.OrdersClient
}

func NewClient(o *common.ClientOptions) *Client {
	deviceClient := devices.NewDevicesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&deviceClient.Client, o.ResourceManagerAuthorizer)

	orderClient := databoxedge.NewOrdersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&orderClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		DeviceClient: &deviceClient,
		OrderClient:  &orderClient,
	}
}

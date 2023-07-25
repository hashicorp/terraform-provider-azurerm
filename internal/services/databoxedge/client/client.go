// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/databoxedge/2022-03-01/devices"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databoxedge/2022-03-01/orders"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	DeviceClient *devices.DevicesClient
	OrdersClient *orders.OrdersClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	deviceClient, err := devices.NewDevicesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Devices client: %+v", err)
	}
	o.Configure(deviceClient.Client, o.Authorizers.ResourceManager)

	ordersClient, err := orders.NewOrdersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Orders client: %+v", err)
	}
	o.Configure(ordersClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		DeviceClient: deviceClient,
		OrdersClient: ordersClient,
	}, nil
}

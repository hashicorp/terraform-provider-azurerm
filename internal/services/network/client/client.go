// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2022-07-01/applicationgateways"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/networkinterfaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/vmsspublicipaddresses"
	network_2023_11_01 "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	*network_2023_11_01.Client

	// TODO 4.0 application gateways should be updated to use 2023-09-01 just prior to releasing 4.0
	ApplicationGatewaysClient *applicationgateways.ApplicationGatewaysClient
	// VMSS Data Source requires the Network Interfaces and VMSSPublicIpAddresses client from `2023-09-01` for the `ListVirtualMachineScaleSetVMNetworkInterfacesComplete` method
	NetworkInterfacesClient     *networkinterfaces.NetworkInterfacesClient
	VMSSPublicIPAddressesClient *vmsspublicipaddresses.VMSSPublicIPAddressesClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	ApplicationGatewaysClient, err := applicationgateways.NewApplicationGatewaysClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Application Gateways Client: %+v", err)
	}
	o.Configure(ApplicationGatewaysClient.Client, o.Authorizers.ResourceManager)

	NetworkInterfacesClient, err := networkinterfaces.NewNetworkInterfacesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Network Interfaces Client: %+v", err)
	}
	o.Configure(NetworkInterfacesClient.Client, o.Authorizers.ResourceManager)

	VMSSPublicIPAddressesClient, err := vmsspublicipaddresses.NewVMSSPublicIPAddressesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building VMSS Public IP Addresses Client: %+v", err)
	}
	o.Configure(VMSSPublicIPAddressesClient.Client, o.Authorizers.ResourceManager)

	client, err := network_2023_11_01.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building clients for Network: %+v", err)
	}

	return &Client{
		ApplicationGatewaysClient:   ApplicationGatewaysClient,
		NetworkInterfacesClient:     NetworkInterfacesClient,
		VMSSPublicIPAddressesClient: VMSSPublicIPAddressesClient,
		Client:                      client,
	}, nil
}

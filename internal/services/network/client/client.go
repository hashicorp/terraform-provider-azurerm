// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/networkinterfaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/vmsspublicipaddresses"
	network_2023_11_01 "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-01-01/bastionhosts"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	*network_2023_11_01.Client

	BastionHostsClient *bastionhosts.BastionHostsClient
	// VMSS Data Source requires the Network Interfaces and VMSSPublicIpAddresses client from `2023-09-01` for the `ListVirtualMachineScaleSetVMNetworkInterfacesComplete` method
	NetworkInterfacesClient     *networkinterfaces.NetworkInterfacesClient
	VMSSPublicIPAddressesClient *vmsspublicipaddresses.VMSSPublicIPAddressesClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	BastionHostsClient, err := bastionhosts.NewBastionHostsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Bastion Client: %+v", err)
	}
	o.Configure(BastionHostsClient.Client, o.Authorizers.ResourceManager)

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
		BastionHostsClient:          BastionHostsClient,
		NetworkInterfacesClient:     NetworkInterfacesClient,
		VMSSPublicIPAddressesClient: VMSSPublicIPAddressesClient,
		Client:                      client,
	}, nil
}

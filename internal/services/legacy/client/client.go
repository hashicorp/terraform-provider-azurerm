// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/tombuildsstuff/kermit/sdk/compute/2023-03-01/compute"
)

type Client struct {
	DisksClient      *compute.DisksClient
	VMClient         *compute.VirtualMachinesClient
	VMScaleSetClient *compute.VirtualMachineScaleSetsClient
}

func NewClient(o *common.ClientOptions) *Client {
	disksClient := compute.NewDisksClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&disksClient.Client, o.ResourceManagerAuthorizer)

	vmScaleSetClient := compute.NewVirtualMachineScaleSetsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&vmScaleSetClient.Client, o.ResourceManagerAuthorizer)

	vmClient := compute.NewVirtualMachinesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&vmClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		DisksClient:      &disksClient,
		VMScaleSetClient: &vmScaleSetClient,
		VMClient:         &vmClient,
	}
}

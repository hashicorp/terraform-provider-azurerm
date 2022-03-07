package client

import (
	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2021-07-01/compute"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
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

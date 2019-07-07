package devtestlabs

import (
	"github.com/Azure/azure-sdk-for-go/services/devtestlabs/mgmt/2016-05-15/dtl"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	LabsClient            dtl.LabsClient
	PoliciesClient        dtl.PoliciesClient
	VirtualMachinesClient dtl.VirtualMachinesClient
	VirtualNetworksClient dtl.VirtualNetworksClient
}

func BuildClient(o *common.ClientOptions) *Client {
	c := Client{}

	c.LabsClient = dtl.NewLabsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.LabsClient.Client, o.ResourceManagerAuthorizer)

	c.PoliciesClient = dtl.NewPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.PoliciesClient.Client, o.ResourceManagerAuthorizer)

	c.VirtualMachinesClient = dtl.NewVirtualMachinesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.VirtualMachinesClient.Client, o.ResourceManagerAuthorizer)

	c.VirtualNetworksClient = dtl.NewVirtualNetworksClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.VirtualNetworksClient.Client, o.ResourceManagerAuthorizer)

	return &c
}

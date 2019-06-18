package devtestlabs

import (
	"github.com/Azure/azure-sdk-for-go/services/devtestlabs/mgmt/2016-05-15/dtl"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/ar"
)

type Client struct {
	LabsClient            dtl.LabsClient
	PoliciesClient        dtl.PoliciesClient
	VirtualMachinesClient dtl.VirtualMachinesClient
	VirtualNetworksClient dtl.VirtualNetworksClient
}

func BuildClient(endpoint, subscriptionId string, o *ar.ClientOptions) *Client {
	c := Client{}

	c.LabsClient = dtl.NewLabsClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.LabsClient.Client, o)

	c.PoliciesClient = dtl.NewPoliciesClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.PoliciesClient.Client, o)

	c.VirtualMachinesClient = dtl.NewVirtualMachinesClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.VirtualMachinesClient.Client, o)

	c.VirtualNetworksClient = dtl.NewVirtualNetworksClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.VirtualNetworksClient.Client, o)

	return &c
}

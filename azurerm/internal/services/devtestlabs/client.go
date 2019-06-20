package devtestlabs

import (
	"github.com/Azure/azure-sdk-for-go/services/devtestlabs/mgmt/2016-05-15/dtl"
	"github.com/Azure/go-autorest/autorest"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	LabsClient            dtl.LabsClient
	PoliciesClient        dtl.PoliciesClient
	VirtualMachinesClient dtl.VirtualMachinesClient
	VirtualNetworksClient dtl.VirtualNetworksClient
}

func BuildClient(endpoint string, authorizer autorest.Authorizer, o *common.ClientOptions) *Client {
	c := Client{}

	c.LabsClient = dtl.NewLabsClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.LabsClient.Client, authorizer)

	c.PoliciesClient = dtl.NewPoliciesClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.PoliciesClient.Client, authorizer)

	c.VirtualMachinesClient = dtl.NewVirtualMachinesClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.VirtualMachinesClient.Client, authorizer)

	c.VirtualNetworksClient = dtl.NewVirtualNetworksClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.VirtualNetworksClient.Client, authorizer)

	return &c
}

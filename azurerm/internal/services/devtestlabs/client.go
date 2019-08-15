package devtestlabs

import (
	"github.com/Azure/azure-sdk-for-go/services/devtestlabs/mgmt/2016-05-15/dtl"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	LabsClient            *dtl.LabsClient
	LabSchedulesClient    *dtl.SchedulesClient
	PoliciesClient        *dtl.PoliciesClient
	VirtualMachinesClient *dtl.VirtualMachinesClient
	VirtualNetworksClient *dtl.VirtualNetworksClient
}

func BuildClient(o *common.ClientOptions) *Client {

	LabsClient := dtl.NewLabsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&LabsClient.Client, o.ResourceManagerAuthorizer)

	PoliciesClient := dtl.NewPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&PoliciesClient.Client, o.ResourceManagerAuthorizer)

	VirtualMachinesClient := dtl.NewVirtualMachinesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&VirtualMachinesClient.Client, o.ResourceManagerAuthorizer)

	VirtualNetworksClient := dtl.NewVirtualNetworksClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&VirtualNetworksClient.Client, o.ResourceManagerAuthorizer)

	LabSchedulesClient := dtl.NewSchedulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&LabSchedulesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		LabsClient:            &LabsClient,
		LabSchedulesClient:    &LabSchedulesClient,
		PoliciesClient:        &PoliciesClient,
		VirtualMachinesClient: &VirtualMachinesClient,
		VirtualNetworksClient: &VirtualNetworksClient,
	}
}

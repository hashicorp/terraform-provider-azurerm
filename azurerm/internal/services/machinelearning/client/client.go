package machinelearning

import (
	"github.com/Azure/azure-sdk-for-go/services/machinelearningservices/mgmt/2019-11-01/machinelearningservices"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	WorkspacesClient             *machinelearningservices.WorkspacesClient
	VirtualMachineSizesClient    *machinelearningservices.VirtualMachineSizesClient
	MachineLearningComputeClient *machinelearningservices.MachineLearningComputeClient
	UsagesClient                 *machinelearningservices.UsagesClient
}

func NewClient(o *common.ClientOptions) *Client {

	WorkspacesClient := machinelearningservices.NewWorkspacesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&WorkspacesClient.Client, o.ResourceManagerAuthorizer)

	VirtualMachineSizesClient := machinelearningservices.NewVirtualMachineSizesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&VirtualMachineSizesClient.Client, o.ResourceManagerAuthorizer)

	MachineLearningComputeClient := machinelearningservices.NewMachineLearningComputeClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&MachineLearningComputeClient.Client, o.ResourceManagerAuthorizer)

	UsagesClient := machinelearningservices.NewUsagesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&UsagesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		WorkspacesClient:             &WorkspacesClient,
		VirtualMachineSizesClient:    &VirtualMachineSizesClient,
		MachineLearningComputeClient: &MachineLearningComputeClient,
		UsagesClient:                 &UsagesClient,
	}
}

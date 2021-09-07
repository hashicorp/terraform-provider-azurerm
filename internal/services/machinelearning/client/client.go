package client

import (
	"github.com/Azure/azure-sdk-for-go/services/machinelearningservices/mgmt/2021-07-01/machinelearningservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	WorkspacesClient             *machinelearningservices.WorkspacesClient
	MachineLearningComputeClient *machinelearningservices.ComputeClient
}

func NewClient(o *common.ClientOptions) *Client {
	WorkspacesClient := machinelearningservices.NewWorkspacesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&WorkspacesClient.Client, o.ResourceManagerAuthorizer)

	MachineLearningComputeClient := machinelearningservices.NewComputeClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&MachineLearningComputeClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		WorkspacesClient:             &WorkspacesClient,
		MachineLearningComputeClient: &MachineLearningComputeClient,
	}
}

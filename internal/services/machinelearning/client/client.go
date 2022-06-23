package client

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/sdk/2021-07-01/machinelearningcomputes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/sdk/2021-07-01/workspaces"
)

type Client struct {
	ComputeClient    *machinelearningcomputes.MachineLearningComputesClient
	WorkspacesClient *workspaces.WorkspacesClient
}

func NewClient(o *common.ClientOptions) *Client {
	ComputeClient := machinelearningcomputes.NewMachineLearningComputesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&ComputeClient.Client, o.ResourceManagerAuthorizer)

	WorkspacesClient := workspaces.NewWorkspacesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&WorkspacesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ComputeClient:    &ComputeClient,
		WorkspacesClient: &WorkspacesClient,
	}
}

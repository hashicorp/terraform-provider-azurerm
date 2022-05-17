package client

import (
	"github.com/Azure/azure-sdk-for-go/services/machinelearningservices/mgmt/2021-07-01/machinelearningservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ComputeClient    *machinelearningservices.ComputeClient
	WorkspacesClient *machinelearningservices.WorkspacesClient
}

func NewClient(o *common.ClientOptions) *Client {
	ComputeClient := machinelearningservices.NewComputeClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ComputeClient.Client, o.ResourceManagerAuthorizer)

	WorkspacesClient := machinelearningservices.NewWorkspacesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&WorkspacesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ComputeClient:    &ComputeClient,
		WorkspacesClient: &WorkspacesClient,
	}
}

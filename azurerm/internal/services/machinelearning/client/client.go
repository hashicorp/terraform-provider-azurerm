package client

import (
	"github.com/Azure/azure-sdk-for-go/services/machinelearningservices/mgmt/2020-04-01/machinelearningservices"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	WorkspacesClient *machinelearningservices.WorkspacesClient
}

func NewClient(o *common.ClientOptions) *Client {
	WorkspacesClient := machinelearningservices.NewWorkspacesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&WorkspacesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		WorkspacesClient: &WorkspacesClient,
	}
}

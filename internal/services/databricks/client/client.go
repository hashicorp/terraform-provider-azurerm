package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2021-04-01-preview/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	WorkspacesClient *workspaces.WorkspacesClient
}

func NewClient(o *common.ClientOptions) *Client {
	WorkspacesClient := workspaces.NewWorkspacesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&WorkspacesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		WorkspacesClient: &WorkspacesClient,
	}
}

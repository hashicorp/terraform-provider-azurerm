package client

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/databricks/sdk/2021-04-01-preview/workspaces"
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

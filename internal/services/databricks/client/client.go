package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2022-04-01-preview/accessconnector"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2023-02-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AccessConnectorClient *accessconnector.AccessConnectorClient
	WorkspacesClient      *workspaces.WorkspacesClient
}

func NewClient(o *common.ClientOptions) *Client {
	AccessConnectorClient := accessconnector.NewAccessConnectorClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&AccessConnectorClient.Client, o.ResourceManagerAuthorizer)
	WorkspacesClient := workspaces.NewWorkspacesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&WorkspacesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AccessConnectorClient: &AccessConnectorClient,
		WorkspacesClient:      &WorkspacesClient,
	}
}

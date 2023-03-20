package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2022-04-01-preview/accessconnector"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2023-02-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AccessConnectorClient *accessconnector.AccessConnectorClient
	WorkspacesClient      *workspaces.WorkspacesClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	accessConnectorClient, err := accessconnector.NewAccessConnectorClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building AccessConnector client: %+v", err)
	}
	o.Configure(accessConnectorClient.Client, o.Authorizers.ResourceManager)

	workspacesClient, err := workspaces.NewWorkspacesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Workspaces client: %+v", err)
	}
	o.Configure(workspacesClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		AccessConnectorClient: accessConnectorClient,
		WorkspacesClient:      workspacesClient,
	}, nil
}

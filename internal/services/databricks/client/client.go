package client

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/databricks/mgmt/2021-04-01-preview/databricks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	WorkspacesClient *databricks.WorkspacesClient
}

func NewClient(o *common.ClientOptions) *Client {
	WorkspacesClient := databricks.NewWorkspacesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&WorkspacesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		WorkspacesClient: &WorkspacesClient,
	}
}

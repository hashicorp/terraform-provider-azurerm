package databricks

import (
	"github.com/Azure/azure-sdk-for-go/services/databricks/mgmt/2018-04-01/databricks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	WorkspacesClient *databricks.WorkspacesClient
}

func BuildClient(o *common.ClientOptions) *Client {

	WorkspacesClient := databricks.NewWorkspacesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&WorkspacesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		WorkspacesClient: &WorkspacesClient,
	}
}

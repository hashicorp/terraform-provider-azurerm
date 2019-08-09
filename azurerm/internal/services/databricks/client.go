package databricks

import (
	"github.com/Azure/azure-sdk-for-go/services/databricks/mgmt/2018-04-01/databricks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	WorkspacesClient databricks.WorkspacesClient
}

func BuildClient(o *common.ClientOptions) *Client {
	c := Client{}

	c.WorkspacesClient = databricks.NewWorkspacesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.WorkspacesClient.Client, o.ResourceManagerAuthorizer)

	return &c
}

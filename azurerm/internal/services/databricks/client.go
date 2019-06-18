package databricks

import (
	"github.com/Azure/azure-sdk-for-go/services/databricks/mgmt/2018-04-01/databricks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/ar"
)

type Client struct {
	WorkspacesClient databricks.WorkspacesClient
}

func BuildClient(endpoint, subscriptionId string, o *ar.ClientOptions) *Client {
	c := Client{}

	c.WorkspacesClient = databricks.NewWorkspacesClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.WorkspacesClient.Client, o)

	return &c
}

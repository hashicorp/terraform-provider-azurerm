package databricks

import (
	"github.com/Azure/azure-sdk-for-go/services/databricks/mgmt/2018-04-01/databricks"
	"github.com/Azure/go-autorest/autorest"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	WorkspacesClient databricks.WorkspacesClient
}

func BuildClient(endpoint string, authorizer autorest.Authorizer, o *common.ClientOptions) *Client {
	c := Client{}

	c.WorkspacesClient = databricks.NewWorkspacesClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.WorkspacesClient.Client, authorizer)

	return &c
}

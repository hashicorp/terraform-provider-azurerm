package databricks

import "github.com/Azure/azure-sdk-for-go/services/databricks/mgmt/2018-04-01/databricks"

type Client struct {
	WorkspacesClient databricks.WorkspacesClient
}

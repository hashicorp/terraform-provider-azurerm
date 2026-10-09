
## `github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/databaseoperations` Documentation

The `databaseoperations` SDK allows for interaction with Azure Resource Manager `sql` (API Version `2025-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/databaseoperations"
```


### Client Initialization

```go
client := databaseoperations.NewDatabaseOperationsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DatabaseOperationsClient.ListByDatabase`

```go
ctx := context.TODO()
id := commonids.NewSqlDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName", "databaseName")

// alternatively `client.ListByDatabase(ctx, id)` can be used to do batched pagination
items, err := client.ListByDatabaseComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

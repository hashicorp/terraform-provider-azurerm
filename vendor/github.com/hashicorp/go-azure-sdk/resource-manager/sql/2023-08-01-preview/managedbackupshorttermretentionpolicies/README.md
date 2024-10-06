
## `github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/managedbackupshorttermretentionpolicies` Documentation

The `managedbackupshorttermretentionpolicies` SDK allows for interaction with Azure Resource Manager `sql` (API Version `2023-08-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/managedbackupshorttermretentionpolicies"
```


### Client Initialization

```go
client := managedbackupshorttermretentionpolicies.NewManagedBackupShortTermRetentionPoliciesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ManagedBackupShortTermRetentionPoliciesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := commonids.NewSqlManagedInstanceDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedInstanceName", "databaseName")

payload := managedbackupshorttermretentionpolicies.ManagedBackupShortTermRetentionPolicy{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ManagedBackupShortTermRetentionPoliciesClient.Get`

```go
ctx := context.TODO()
id := commonids.NewSqlManagedInstanceDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedInstanceName", "databaseName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedBackupShortTermRetentionPoliciesClient.ListByDatabase`

```go
ctx := context.TODO()
id := commonids.NewSqlManagedInstanceDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedInstanceName", "databaseName")

// alternatively `client.ListByDatabase(ctx, id)` can be used to do batched pagination
items, err := client.ListByDatabaseComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ManagedBackupShortTermRetentionPoliciesClient.Update`

```go
ctx := context.TODO()
id := commonids.NewSqlManagedInstanceDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedInstanceName", "databaseName")

payload := managedbackupshorttermretentionpolicies.ManagedBackupShortTermRetentionPolicy{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

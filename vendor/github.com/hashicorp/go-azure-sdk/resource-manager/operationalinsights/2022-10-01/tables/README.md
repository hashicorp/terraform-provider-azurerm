
## `github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/tables` Documentation

The `tables` SDK allows for interaction with Azure Resource Manager `operationalinsights` (API Version `2022-10-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/tables"
```


### Client Initialization

```go
client := tables.NewTablesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `TablesClient.CancelSearch`

```go
ctx := context.TODO()
id := tables.NewTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "tableName")

read, err := client.CancelSearch(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TablesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := tables.NewTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "tableName")

payload := tables.Table{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `TablesClient.Delete`

```go
ctx := context.TODO()
id := tables.NewTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "tableName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `TablesClient.Get`

```go
ctx := context.TODO()
id := tables.NewTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "tableName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TablesClient.ListByWorkspace`

```go
ctx := context.TODO()
id := tables.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName")

read, err := client.ListByWorkspace(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TablesClient.Migrate`

```go
ctx := context.TODO()
id := tables.NewTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "tableName")

read, err := client.Migrate(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TablesClient.Update`

```go
ctx := context.TODO()
id := tables.NewTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "tableName")

payload := tables.Table{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

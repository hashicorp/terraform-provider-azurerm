
## `github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2022-01-01/databases` Documentation

The `databases` SDK allows for interaction with Azure Resource Manager `mysql` (API Version `2022-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2022-01-01/databases"
```


### Client Initialization

```go
client := databases.NewDatabasesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DatabasesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := databases.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerName", "databaseName")

payload := databases.Database{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DatabasesClient.Delete`

```go
ctx := context.TODO()
id := databases.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerName", "databaseName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `DatabasesClient.Get`

```go
ctx := context.TODO()
id := databases.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerName", "databaseName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DatabasesClient.ListByServer`

```go
ctx := context.TODO()
id := databases.NewFlexibleServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerName")

// alternatively `client.ListByServer(ctx, id)` can be used to do batched pagination
items, err := client.ListByServerComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


## `github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2026-06-01/databaseeditions` Documentation

The `databaseeditions` SDK allows for interaction with Azure Resource Manager `oracledatabase` (API Version `2026-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2026-06-01/databaseeditions"
```


### Client Initialization

```go
client := databaseeditions.NewDatabaseEditionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DatabaseEditionsClient.Get`

```go
ctx := context.TODO()
id := databaseeditions.NewDatabaseEditionID("12345678-1234-9876-4563-123456789012", "locationName", "databaseEditionName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DatabaseEditionsClient.ListByLocation`

```go
ctx := context.TODO()
id := databaseeditions.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

// alternatively `client.ListByLocation(ctx, id)` can be used to do batched pagination
items, err := client.ListByLocationComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


## `github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/bookshelfprivateendpointconnections` Documentation

The `bookshelfprivateendpointconnections` SDK allows for interaction with Azure Resource Manager `discovery` (API Version `2026-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/bookshelfprivateendpointconnections"
```


### Client Initialization

```go
client := bookshelfprivateendpointconnections.NewBookshelfPrivateEndpointConnectionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `BookshelfPrivateEndpointConnectionsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := bookshelfprivateendpointconnections.NewBookshelfPrivateEndpointConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "bookshelfName", "privateEndpointConnectionName")

payload := bookshelfprivateendpointconnections.BookshelfPrivateEndpointConnection{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `BookshelfPrivateEndpointConnectionsClient.Delete`

```go
ctx := context.TODO()
id := bookshelfprivateendpointconnections.NewBookshelfPrivateEndpointConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "bookshelfName", "privateEndpointConnectionName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `BookshelfPrivateEndpointConnectionsClient.Get`

```go
ctx := context.TODO()
id := bookshelfprivateendpointconnections.NewBookshelfPrivateEndpointConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "bookshelfName", "privateEndpointConnectionName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BookshelfPrivateEndpointConnectionsClient.ListByBookshelf`

```go
ctx := context.TODO()
id := bookshelfprivateendpointconnections.NewBookshelfID("12345678-1234-9876-4563-123456789012", "example-resource-group", "bookshelfName")

// alternatively `client.ListByBookshelf(ctx, id)` can be used to do batched pagination
items, err := client.ListByBookshelfComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

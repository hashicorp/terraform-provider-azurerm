
## `github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/bookshelfprivatelinkresources` Documentation

The `bookshelfprivatelinkresources` SDK allows for interaction with Azure Resource Manager `discovery` (API Version `2026-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/bookshelfprivatelinkresources"
```


### Client Initialization

```go
client := bookshelfprivatelinkresources.NewBookshelfPrivateLinkResourcesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `BookshelfPrivateLinkResourcesClient.Get`

```go
ctx := context.TODO()
id := bookshelfprivatelinkresources.NewBookshelfPrivateLinkResourceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "bookshelfName", "privateLinkResourceName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BookshelfPrivateLinkResourcesClient.ListByBookshelf`

```go
ctx := context.TODO()
id := bookshelfprivatelinkresources.NewBookshelfID("12345678-1234-9876-4563-123456789012", "example-resource-group", "bookshelfName")

// alternatively `client.ListByBookshelf(ctx, id)` can be used to do batched pagination
items, err := client.ListByBookshelfComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

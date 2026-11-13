
## `github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/bookshelves` Documentation

The `bookshelves` SDK allows for interaction with Azure Resource Manager `discovery` (API Version `2026-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/bookshelves"
```


### Client Initialization

```go
client := bookshelves.NewBookshelvesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `BookshelvesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := bookshelves.NewBookshelfID("12345678-1234-9876-4563-123456789012", "example-resource-group", "bookshelfName")

payload := bookshelves.Bookshelf{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `BookshelvesClient.Delete`

```go
ctx := context.TODO()
id := bookshelves.NewBookshelfID("12345678-1234-9876-4563-123456789012", "example-resource-group", "bookshelfName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `BookshelvesClient.Get`

```go
ctx := context.TODO()
id := bookshelves.NewBookshelfID("12345678-1234-9876-4563-123456789012", "example-resource-group", "bookshelfName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BookshelvesClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `BookshelvesClient.ListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `BookshelvesClient.Update`

```go
ctx := context.TODO()
id := bookshelves.NewBookshelfID("12345678-1234-9876-4563-123456789012", "example-resource-group", "bookshelfName")

payload := bookshelves.BookshelfUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

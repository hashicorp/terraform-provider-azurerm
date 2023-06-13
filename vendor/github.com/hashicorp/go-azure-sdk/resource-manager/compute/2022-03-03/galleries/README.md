
## `github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleries` Documentation

The `galleries` SDK allows for interaction with the Azure Resource Manager Service `compute` (API Version `2022-03-03`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleries"
```


### Client Initialization

```go
client := galleries.NewGalleriesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `GalleriesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := galleries.NewGalleryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "galleryValue")

payload := galleries.Gallery{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `GalleriesClient.Delete`

```go
ctx := context.TODO()
id := galleries.NewGalleryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "galleryValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `GalleriesClient.Get`

```go
ctx := context.TODO()
id := galleries.NewGalleryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "galleryValue")

read, err := client.Get(ctx, id, galleries.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GalleriesClient.List`

```go
ctx := context.TODO()
id := galleries.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `GalleriesClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := galleries.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `GalleriesClient.Update`

```go
ctx := context.TODO()
id := galleries.NewGalleryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "galleryValue")

payload := galleries.GalleryUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

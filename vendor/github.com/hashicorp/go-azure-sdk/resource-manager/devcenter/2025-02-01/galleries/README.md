
## `github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/galleries` Documentation

The `galleries` SDK allows for interaction with Azure Resource Manager `devcenter` (API Version `2025-02-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/galleries"
```


### Client Initialization

```go
client := galleries.NewGalleriesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `GalleriesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := galleries.NewGalleryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "devCenterName", "galleryName")

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
id := galleries.NewGalleryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "devCenterName", "galleryName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `GalleriesClient.Get`

```go
ctx := context.TODO()
id := galleries.NewGalleryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "devCenterName", "galleryName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GalleriesClient.ListByDevCenter`

```go
ctx := context.TODO()
id := galleries.NewDevCenterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "devCenterName")

// alternatively `client.ListByDevCenter(ctx, id, galleries.DefaultListByDevCenterOperationOptions())` can be used to do batched pagination
items, err := client.ListByDevCenterComplete(ctx, id, galleries.DefaultListByDevCenterOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


## `github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/images` Documentation

The `images` SDK allows for interaction with Azure Resource Manager `devcenter` (API Version `2025-02-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/images"
```


### Client Initialization

```go
client := images.NewImagesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ImagesClient.Get`

```go
ctx := context.TODO()
id := images.NewGalleryImageID("12345678-1234-9876-4563-123456789012", "example-resource-group", "devCenterName", "galleryName", "imageName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ImagesClient.GetByProject`

```go
ctx := context.TODO()
id := images.NewImageID("12345678-1234-9876-4563-123456789012", "example-resource-group", "projectName", "imageName")

read, err := client.GetByProject(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ImagesClient.ListByDevCenter`

```go
ctx := context.TODO()
id := images.NewDevCenterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "devCenterName")

// alternatively `client.ListByDevCenter(ctx, id, images.DefaultListByDevCenterOperationOptions())` can be used to do batched pagination
items, err := client.ListByDevCenterComplete(ctx, id, images.DefaultListByDevCenterOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ImagesClient.ListByGallery`

```go
ctx := context.TODO()
id := images.NewGalleryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "devCenterName", "galleryName")

// alternatively `client.ListByGallery(ctx, id, images.DefaultListByGalleryOperationOptions())` can be used to do batched pagination
items, err := client.ListByGalleryComplete(ctx, id, images.DefaultListByGalleryOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ImagesClient.ListByProject`

```go
ctx := context.TODO()
id := images.NewProjectID("12345678-1234-9876-4563-123456789012", "example-resource-group", "projectName")

// alternatively `client.ListByProject(ctx, id)` can be used to do batched pagination
items, err := client.ListByProjectComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


## `github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryimages` Documentation

The `galleryimages` SDK allows for interaction with Azure Resource Manager `compute` (API Version `2022-03-03`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryimages"
```


### Client Initialization

```go
client := galleryimages.NewGalleryImagesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `GalleryImagesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := galleryimages.NewGalleryImageID("12345678-1234-9876-4563-123456789012", "example-resource-group", "galleryName", "imageName")

payload := galleryimages.GalleryImage{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `GalleryImagesClient.Delete`

```go
ctx := context.TODO()
id := galleryimages.NewGalleryImageID("12345678-1234-9876-4563-123456789012", "example-resource-group", "galleryName", "imageName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `GalleryImagesClient.Get`

```go
ctx := context.TODO()
id := galleryimages.NewGalleryImageID("12345678-1234-9876-4563-123456789012", "example-resource-group", "galleryName", "imageName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GalleryImagesClient.ListByGallery`

```go
ctx := context.TODO()
id := commonids.NewSharedImageGalleryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "galleryName")

// alternatively `client.ListByGallery(ctx, id)` can be used to do batched pagination
items, err := client.ListByGalleryComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `GalleryImagesClient.Update`

```go
ctx := context.TODO()
id := galleryimages.NewGalleryImageID("12345678-1234-9876-4563-123456789012", "example-resource-group", "galleryName", "imageName")

payload := galleryimages.GalleryImageUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

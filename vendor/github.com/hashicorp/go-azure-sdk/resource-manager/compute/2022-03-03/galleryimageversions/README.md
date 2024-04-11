
## `github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryimageversions` Documentation

The `galleryimageversions` SDK allows for interaction with the Azure Resource Manager Service `compute` (API Version `2022-03-03`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryimageversions"
```


### Client Initialization

```go
client := galleryimageversions.NewGalleryImageVersionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `GalleryImageVersionsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := galleryimageversions.NewImageVersionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "galleryValue", "imageValue", "versionValue")

payload := galleryimageversions.GalleryImageVersion{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `GalleryImageVersionsClient.Delete`

```go
ctx := context.TODO()
id := galleryimageversions.NewImageVersionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "galleryValue", "imageValue", "versionValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `GalleryImageVersionsClient.Get`

```go
ctx := context.TODO()
id := galleryimageversions.NewImageVersionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "galleryValue", "imageValue", "versionValue")

read, err := client.Get(ctx, id, galleryimageversions.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GalleryImageVersionsClient.ListByGalleryImage`

```go
ctx := context.TODO()
id := galleryimageversions.NewGalleryImageID("12345678-1234-9876-4563-123456789012", "example-resource-group", "galleryValue", "imageValue")

// alternatively `client.ListByGalleryImage(ctx, id)` can be used to do batched pagination
items, err := client.ListByGalleryImageComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `GalleryImageVersionsClient.Update`

```go
ctx := context.TODO()
id := galleryimageversions.NewImageVersionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "galleryValue", "imageValue", "versionValue")

payload := galleryimageversions.GalleryImageVersionUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

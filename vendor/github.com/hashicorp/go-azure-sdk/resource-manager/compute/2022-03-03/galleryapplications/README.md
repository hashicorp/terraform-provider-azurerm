
## `github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryapplications` Documentation

The `galleryapplications` SDK allows for interaction with the Azure Resource Manager Service `compute` (API Version `2022-03-03`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryapplications"
```


### Client Initialization

```go
client := galleryapplications.NewGalleryApplicationsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `GalleryApplicationsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := galleryapplications.NewApplicationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "galleryValue", "applicationValue")

payload := galleryapplications.GalleryApplication{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `GalleryApplicationsClient.Delete`

```go
ctx := context.TODO()
id := galleryapplications.NewApplicationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "galleryValue", "applicationValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `GalleryApplicationsClient.Get`

```go
ctx := context.TODO()
id := galleryapplications.NewApplicationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "galleryValue", "applicationValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GalleryApplicationsClient.ListByGallery`

```go
ctx := context.TODO()
id := galleryapplications.NewGalleryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "galleryValue")

// alternatively `client.ListByGallery(ctx, id)` can be used to do batched pagination
items, err := client.ListByGalleryComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `GalleryApplicationsClient.Update`

```go
ctx := context.TODO()
id := galleryapplications.NewApplicationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "galleryValue", "applicationValue")

payload := galleryapplications.GalleryApplicationUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


## `github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryapplicationversions` Documentation

The `galleryapplicationversions` SDK allows for interaction with Azure Resource Manager `compute` (API Version `2022-03-03`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryapplicationversions"
```


### Client Initialization

```go
client := galleryapplicationversions.NewGalleryApplicationVersionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `GalleryApplicationVersionsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := galleryapplicationversions.NewApplicationVersionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "galleryName", "applicationName", "versionName")

payload := galleryapplicationversions.GalleryApplicationVersion{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `GalleryApplicationVersionsClient.Delete`

```go
ctx := context.TODO()
id := galleryapplicationversions.NewApplicationVersionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "galleryName", "applicationName", "versionName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `GalleryApplicationVersionsClient.Get`

```go
ctx := context.TODO()
id := galleryapplicationversions.NewApplicationVersionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "galleryName", "applicationName", "versionName")

read, err := client.Get(ctx, id, galleryapplicationversions.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GalleryApplicationVersionsClient.ListByGalleryApplication`

```go
ctx := context.TODO()
id := galleryapplicationversions.NewApplicationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "galleryName", "applicationName")

// alternatively `client.ListByGalleryApplication(ctx, id)` can be used to do batched pagination
items, err := client.ListByGalleryApplicationComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `GalleryApplicationVersionsClient.Update`

```go
ctx := context.TODO()
id := galleryapplicationversions.NewApplicationVersionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "galleryName", "applicationName", "versionName")

payload := galleryapplicationversions.GalleryApplicationVersionUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

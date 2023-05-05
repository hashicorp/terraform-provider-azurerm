
## `github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/images` Documentation

The `images` SDK allows for interaction with the Azure Resource Manager Service `compute` (API Version `2022-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/images"
```


### Client Initialization

```go
client := images.NewImagesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ImagesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := images.NewImageID("12345678-1234-9876-4563-123456789012", "example-resource-group", "imageValue")

payload := images.Image{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ImagesClient.Delete`

```go
ctx := context.TODO()
id := images.NewImageID("12345678-1234-9876-4563-123456789012", "example-resource-group", "imageValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ImagesClient.Get`

```go
ctx := context.TODO()
id := images.NewImageID("12345678-1234-9876-4563-123456789012", "example-resource-group", "imageValue")

read, err := client.Get(ctx, id, images.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ImagesClient.List`

```go
ctx := context.TODO()
id := images.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ImagesClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := images.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ImagesClient.Update`

```go
ctx := context.TODO()
id := images.NewImageID("12345678-1234-9876-4563-123456789012", "example-resource-group", "imageValue")

payload := images.ImageUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

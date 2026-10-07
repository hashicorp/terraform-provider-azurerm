
## `github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/marketplacegalleryimages` Documentation

The `marketplacegalleryimages` SDK allows for interaction with Azure Resource Manager `azurestackhci` (API Version `2024-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/marketplacegalleryimages"
```


### Client Initialization

```go
client := marketplacegalleryimages.NewMarketplaceGalleryImagesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `MarketplaceGalleryImagesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := marketplacegalleryimages.NewMarketplaceGalleryImageID("12345678-1234-9876-4563-123456789012", "example-resource-group", "marketplaceGalleryImageName")

payload := marketplacegalleryimages.MarketplaceGalleryImages{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `MarketplaceGalleryImagesClient.Delete`

```go
ctx := context.TODO()
id := marketplacegalleryimages.NewMarketplaceGalleryImageID("12345678-1234-9876-4563-123456789012", "example-resource-group", "marketplaceGalleryImageName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `MarketplaceGalleryImagesClient.Get`

```go
ctx := context.TODO()
id := marketplacegalleryimages.NewMarketplaceGalleryImageID("12345678-1234-9876-4563-123456789012", "example-resource-group", "marketplaceGalleryImageName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `MarketplaceGalleryImagesClient.List`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `MarketplaceGalleryImagesClient.ListAll`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListAll(ctx, id)` can be used to do batched pagination
items, err := client.ListAllComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `MarketplaceGalleryImagesClient.Update`

```go
ctx := context.TODO()
id := marketplacegalleryimages.NewMarketplaceGalleryImageID("12345678-1234-9876-4563-123456789012", "example-resource-group", "marketplaceGalleryImageName")

payload := marketplacegalleryimages.MarketplaceGalleryImagesUpdateRequest{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

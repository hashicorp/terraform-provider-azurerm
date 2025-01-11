
## `github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/skuses` Documentation

The `skuses` SDK allows for interaction with Azure Resource Manager `azurestackhci` (API Version `2024-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/skuses"
```


### Client Initialization

```go
client := skuses.NewSkusesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SkusesClient.SkusGet`

```go
ctx := context.TODO()
id := skuses.NewSkuID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "publisherName", "offerName", "skuName")

read, err := client.SkusGet(ctx, id, skuses.DefaultSkusGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SkusesClient.SkusListByOffer`

```go
ctx := context.TODO()
id := skuses.NewOfferID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "publisherName", "offerName")

// alternatively `client.SkusListByOffer(ctx, id, skuses.DefaultSkusListByOfferOperationOptions())` can be used to do batched pagination
items, err := client.SkusListByOfferComplete(ctx, id, skuses.DefaultSkusListByOfferOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

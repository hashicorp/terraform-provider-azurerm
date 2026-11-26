
## `github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/skus` Documentation

The `skus` SDK allows for interaction with Azure Resource Manager `azurestackhci` (API Version `2024-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/skus"
```


### Client Initialization

```go
client := skus.NewSkusClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SkusClient.Get`

```go
ctx := context.TODO()
id := skus.NewSkuID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "publisherName", "offerName", "skuName")

read, err := client.Get(ctx, id, skus.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SkusClient.ListByOffer`

```go
ctx := context.TODO()
id := skus.NewOfferID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "publisherName", "offerName")

// alternatively `client.ListByOffer(ctx, id, skus.DefaultListByOfferOperationOptions())` can be used to do batched pagination
items, err := client.ListByOfferComplete(ctx, id, skus.DefaultListByOfferOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

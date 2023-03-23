
## `github.com/hashicorp/go-azure-sdk/resource-manager/compute/2021-07-01/skus` Documentation

The `skus` SDK allows for interaction with the Azure Resource Manager Service `compute` (API Version `2021-07-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/compute/2021-07-01/skus"
```


### Client Initialization

```go
client := skus.NewSkusClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SkusClient.ResourceSkusList`

```go
ctx := context.TODO()
id := skus.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ResourceSkusList(ctx, id, skus.DefaultResourceSkusListOperationOptions())` can be used to do batched pagination
items, err := client.ResourceSkusListComplete(ctx, id, skus.DefaultResourceSkusListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

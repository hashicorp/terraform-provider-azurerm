
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/networksecurityperimeterservicetags` Documentation

The `networksecurityperimeterservicetags` SDK allows for interaction with Azure Resource Manager `network` (API Version `2025-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/networksecurityperimeterservicetags"
```


### Client Initialization

```go
client := networksecurityperimeterservicetags.NewNetworkSecurityPerimeterServiceTagsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `NetworkSecurityPerimeterServiceTagsClient.List`

```go
ctx := context.TODO()
id := networksecurityperimeterservicetags.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

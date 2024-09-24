
## `github.com/hashicorp/go-azure-sdk/resource-manager/devopsinfrastructure/2024-04-04-preview/resourcedetails` Documentation

The `resourcedetails` SDK allows for interaction with Azure Resource Manager `devopsinfrastructure` (API Version `2024-04-04-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/devopsinfrastructure/2024-04-04-preview/resourcedetails"
```


### Client Initialization

```go
client := resourcedetails.NewResourceDetailsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ResourceDetailsClient.ListByPool`

```go
ctx := context.TODO()
id := resourcedetails.NewPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "poolName")

// alternatively `client.ListByPool(ctx, id)` can be used to do batched pagination
items, err := client.ListByPoolComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

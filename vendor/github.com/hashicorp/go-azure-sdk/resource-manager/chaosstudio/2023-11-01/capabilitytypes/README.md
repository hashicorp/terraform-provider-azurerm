
## `github.com/hashicorp/go-azure-sdk/resource-manager/chaosstudio/2023-11-01/capabilitytypes` Documentation

The `capabilitytypes` SDK allows for interaction with Azure Resource Manager `chaosstudio` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/chaosstudio/2023-11-01/capabilitytypes"
```


### Client Initialization

```go
client := capabilitytypes.NewCapabilityTypesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `CapabilityTypesClient.List`

```go
ctx := context.TODO()
id := capabilitytypes.NewTargetTypeID("12345678-1234-9876-4563-123456789012", "locationName", "targetTypeName")

// alternatively `client.List(ctx, id, capabilitytypes.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, capabilitytypes.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

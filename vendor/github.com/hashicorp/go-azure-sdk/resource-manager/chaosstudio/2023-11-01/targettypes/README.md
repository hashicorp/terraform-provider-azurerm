
## `github.com/hashicorp/go-azure-sdk/resource-manager/chaosstudio/2023-11-01/targettypes` Documentation

The `targettypes` SDK allows for interaction with Azure Resource Manager `chaosstudio` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/chaosstudio/2023-11-01/targettypes"
```


### Client Initialization

```go
client := targettypes.NewTargetTypesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `TargetTypesClient.Get`

```go
ctx := context.TODO()
id := targettypes.NewTargetTypeID("12345678-1234-9876-4563-123456789012", "locationName", "targetTypeName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TargetTypesClient.List`

```go
ctx := context.TODO()
id := targettypes.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

// alternatively `client.List(ctx, id, targettypes.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, targettypes.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

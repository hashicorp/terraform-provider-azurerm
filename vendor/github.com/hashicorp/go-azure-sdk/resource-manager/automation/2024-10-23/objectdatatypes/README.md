
## `github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/objectdatatypes` Documentation

The `objectdatatypes` SDK allows for interaction with Azure Resource Manager `automation` (API Version `2024-10-23`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/objectdatatypes"
```


### Client Initialization

```go
client := objectdatatypes.NewObjectDataTypesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ObjectDataTypesClient.ListFieldsByModuleAndType`

```go
ctx := context.TODO()
id := objectdatatypes.NewModuleObjectDataTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "moduleName", "objectDataTypeName")

// alternatively `client.ListFieldsByModuleAndType(ctx, id)` can be used to do batched pagination
items, err := client.ListFieldsByModuleAndTypeComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ObjectDataTypesClient.ListFieldsByType`

```go
ctx := context.TODO()
id := objectdatatypes.NewObjectDataTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "objectDataTypeName")

// alternatively `client.ListFieldsByType(ctx, id)` can be used to do batched pagination
items, err := client.ListFieldsByTypeComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

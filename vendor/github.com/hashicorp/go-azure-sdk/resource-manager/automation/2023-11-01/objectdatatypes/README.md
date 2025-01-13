
## `github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/objectdatatypes` Documentation

The `objectdatatypes` SDK allows for interaction with Azure Resource Manager `automation` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/objectdatatypes"
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

read, err := client.ListFieldsByModuleAndType(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ObjectDataTypesClient.ListFieldsByType`

```go
ctx := context.TODO()
id := objectdatatypes.NewObjectDataTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "objectDataTypeName")

read, err := client.ListFieldsByType(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

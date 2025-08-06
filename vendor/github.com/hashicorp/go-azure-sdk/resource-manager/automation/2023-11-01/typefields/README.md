
## `github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/typefields` Documentation

The `typefields` SDK allows for interaction with Azure Resource Manager `automation` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/typefields"
```


### Client Initialization

```go
client := typefields.NewTypeFieldsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `TypeFieldsClient.FieldsListByType`

```go
ctx := context.TODO()
id := typefields.NewTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "automationAccountName", "moduleName", "typeName")

read, err := client.FieldsListByType(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

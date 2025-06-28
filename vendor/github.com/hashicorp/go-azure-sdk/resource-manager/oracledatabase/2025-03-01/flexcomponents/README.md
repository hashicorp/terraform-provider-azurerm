
## `github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-03-01/flexcomponents` Documentation

The `flexcomponents` SDK allows for interaction with Azure Resource Manager `oracledatabase` (API Version `2025-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-03-01/flexcomponents"
```


### Client Initialization

```go
client := flexcomponents.NewFlexComponentsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `FlexComponentsClient.Get`

```go
ctx := context.TODO()
id := flexcomponents.NewFlexComponentID("12345678-1234-9876-4563-123456789012", "locationName", "flexComponentName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FlexComponentsClient.ListByParent`

```go
ctx := context.TODO()
id := flexcomponents.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

// alternatively `client.ListByParent(ctx, id, flexcomponents.DefaultListByParentOperationOptions())` can be used to do batched pagination
items, err := client.ListByParentComplete(ctx, id, flexcomponents.DefaultListByParentOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

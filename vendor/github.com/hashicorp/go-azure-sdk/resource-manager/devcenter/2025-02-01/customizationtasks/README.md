
## `github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/customizationtasks` Documentation

The `customizationtasks` SDK allows for interaction with Azure Resource Manager `devcenter` (API Version `2025-02-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/customizationtasks"
```


### Client Initialization

```go
client := customizationtasks.NewCustomizationTasksClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `CustomizationTasksClient.Get`

```go
ctx := context.TODO()
id := customizationtasks.NewTaskID("12345678-1234-9876-4563-123456789012", "example-resource-group", "devCenterName", "catalogName", "taskName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CustomizationTasksClient.GetErrorDetails`

```go
ctx := context.TODO()
id := customizationtasks.NewTaskID("12345678-1234-9876-4563-123456789012", "example-resource-group", "devCenterName", "catalogName", "taskName")

read, err := client.GetErrorDetails(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CustomizationTasksClient.ListByCatalog`

```go
ctx := context.TODO()
id := customizationtasks.NewDevCenterCatalogID("12345678-1234-9876-4563-123456789012", "example-resource-group", "devCenterName", "catalogName")

// alternatively `client.ListByCatalog(ctx, id)` can be used to do batched pagination
items, err := client.ListByCatalogComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

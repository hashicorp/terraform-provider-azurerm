
## `github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2023-04-01/datacontainer` Documentation

The `datacontainer` SDK allows for interaction with the Azure Resource Manager Service `machinelearningservices` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2023-04-01/datacontainer"
```


### Client Initialization

```go
client := datacontainer.NewDataContainerClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DataContainerClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := datacontainer.NewWorkspaceDataID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "dataValue")

payload := datacontainer.DataContainerResource{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DataContainerClient.Delete`

```go
ctx := context.TODO()
id := datacontainer.NewWorkspaceDataID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "dataValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DataContainerClient.Get`

```go
ctx := context.TODO()
id := datacontainer.NewWorkspaceDataID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "dataValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DataContainerClient.List`

```go
ctx := context.TODO()
id := datacontainer.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue")

// alternatively `client.List(ctx, id, datacontainer.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, datacontainer.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

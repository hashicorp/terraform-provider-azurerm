
## `github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2023-04-01/batchendpoint` Documentation

The `batchendpoint` SDK allows for interaction with the Azure Resource Manager Service `machinelearningservices` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2023-04-01/batchendpoint"
```


### Client Initialization

```go
client := batchendpoint.NewBatchEndpointClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `BatchEndpointClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := batchendpoint.NewBatchEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "batchEndpointValue")

payload := batchendpoint.BatchEndpointTrackedResource{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `BatchEndpointClient.Delete`

```go
ctx := context.TODO()
id := batchendpoint.NewBatchEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "batchEndpointValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `BatchEndpointClient.Get`

```go
ctx := context.TODO()
id := batchendpoint.NewBatchEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "batchEndpointValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BatchEndpointClient.List`

```go
ctx := context.TODO()
id := batchendpoint.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue")

// alternatively `client.List(ctx, id, batchendpoint.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, batchendpoint.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `BatchEndpointClient.ListKeys`

```go
ctx := context.TODO()
id := batchendpoint.NewBatchEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "batchEndpointValue")

read, err := client.ListKeys(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BatchEndpointClient.Update`

```go
ctx := context.TODO()
id := batchendpoint.NewBatchEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "batchEndpointValue")

payload := batchendpoint.PartialMinimalTrackedResourceWithIdentity{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

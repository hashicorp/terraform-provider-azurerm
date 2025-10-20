
## `github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/dataflowendpoint` Documentation

The `dataflowendpoint` SDK allows for interaction with Azure Resource Manager `iotoperations` (API Version `2024-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/dataflowendpoint"
```


### Client Initialization

```go
client := dataflowendpoint.NewDataflowEndpointClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DataflowEndpointClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := dataflowendpoint.NewDataflowEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "instanceName", "dataflowEndpointName")

payload := dataflowendpoint.DataflowEndpointResource{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DataflowEndpointClient.Delete`

```go
ctx := context.TODO()
id := dataflowendpoint.NewDataflowEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "instanceName", "dataflowEndpointName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `DataflowEndpointClient.Get`

```go
ctx := context.TODO()
id := dataflowendpoint.NewDataflowEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "instanceName", "dataflowEndpointName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DataflowEndpointClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := dataflowendpoint.NewInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "instanceName")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

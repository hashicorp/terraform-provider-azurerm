
## `github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/dataflow` Documentation

The `dataflow` SDK allows for interaction with Azure Resource Manager `iotoperations` (API Version `2024-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/dataflow"
```


### Client Initialization

```go
client := dataflow.NewDataflowClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DataflowClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := dataflow.NewDataflowID("12345678-1234-9876-4563-123456789012", "example-resource-group", "instanceName", "dataflowProfileName", "dataflowName")

payload := dataflow.DataflowResource{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DataflowClient.Delete`

```go
ctx := context.TODO()
id := dataflow.NewDataflowID("12345678-1234-9876-4563-123456789012", "example-resource-group", "instanceName", "dataflowProfileName", "dataflowName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `DataflowClient.Get`

```go
ctx := context.TODO()
id := dataflow.NewDataflowID("12345678-1234-9876-4563-123456789012", "example-resource-group", "instanceName", "dataflowProfileName", "dataflowName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DataflowClient.ListByProfileResource`

```go
ctx := context.TODO()
id := dataflow.NewDataflowProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "instanceName", "dataflowProfileName")

// alternatively `client.ListByProfileResource(ctx, id)` can be used to do batched pagination
items, err := client.ListByProfileResourceComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

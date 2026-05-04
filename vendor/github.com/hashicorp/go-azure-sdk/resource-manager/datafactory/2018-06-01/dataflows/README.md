
## `github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/dataflows` Documentation

The `dataflows` SDK allows for interaction with Azure Resource Manager `datafactory` (API Version `2018-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/dataflows"
```


### Client Initialization

```go
client := dataflows.NewDataFlowsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DataFlowsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := dataflows.NewDataflowID("12345678-1234-9876-4563-123456789012", "example-resource-group", "factoryName", "dataflowName")

payload := dataflows.DataFlowResource{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, dataflows.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DataFlowsClient.Delete`

```go
ctx := context.TODO()
id := dataflows.NewDataflowID("12345678-1234-9876-4563-123456789012", "example-resource-group", "factoryName", "dataflowName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DataFlowsClient.Get`

```go
ctx := context.TODO()
id := dataflows.NewDataflowID("12345678-1234-9876-4563-123456789012", "example-resource-group", "factoryName", "dataflowName")

read, err := client.Get(ctx, id, dataflows.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DataFlowsClient.ListByFactory`

```go
ctx := context.TODO()
id := dataflows.NewFactoryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "factoryName")

// alternatively `client.ListByFactory(ctx, id)` can be used to do batched pagination
items, err := client.ListByFactoryComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

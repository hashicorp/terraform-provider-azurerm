
## `github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/pipelines` Documentation

The `pipelines` SDK allows for interaction with Azure Resource Manager `datafactory` (API Version `2018-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/pipelines"
```


### Client Initialization

```go
client := pipelines.NewPipelinesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PipelinesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := pipelines.NewPipelineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "factoryName", "pipelineName")

payload := pipelines.PipelineResource{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, pipelines.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PipelinesClient.CreateRun`

```go
ctx := context.TODO()
id := pipelines.NewPipelineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "factoryName", "pipelineName")
var payload map[string]interface{}

read, err := client.CreateRun(ctx, id, payload, pipelines.DefaultCreateRunOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PipelinesClient.Delete`

```go
ctx := context.TODO()
id := pipelines.NewPipelineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "factoryName", "pipelineName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PipelinesClient.Get`

```go
ctx := context.TODO()
id := pipelines.NewPipelineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "factoryName", "pipelineName")

read, err := client.Get(ctx, id, pipelines.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PipelinesClient.ListByFactory`

```go
ctx := context.TODO()
id := pipelines.NewFactoryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "factoryName")

// alternatively `client.ListByFactory(ctx, id)` can be used to do batched pagination
items, err := client.ListByFactoryComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

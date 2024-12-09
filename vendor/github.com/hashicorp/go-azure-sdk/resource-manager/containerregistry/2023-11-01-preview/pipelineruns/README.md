
## `github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-11-01-preview/pipelineruns` Documentation

The `pipelineruns` SDK allows for interaction with Azure Resource Manager `containerregistry` (API Version `2023-11-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-11-01-preview/pipelineruns"
```


### Client Initialization

```go
client := pipelineruns.NewPipelineRunsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PipelineRunsClient.Create`

```go
ctx := context.TODO()
id := pipelineruns.NewPipelineRunID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName", "pipelineRunName")

payload := pipelineruns.PipelineRun{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `PipelineRunsClient.Delete`

```go
ctx := context.TODO()
id := pipelineruns.NewPipelineRunID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName", "pipelineRunName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `PipelineRunsClient.Get`

```go
ctx := context.TODO()
id := pipelineruns.NewPipelineRunID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName", "pipelineRunName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PipelineRunsClient.List`

```go
ctx := context.TODO()
id := pipelineruns.NewRegistryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


## `github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-06-01-preview/exportpipelines` Documentation

The `exportpipelines` SDK allows for interaction with Azure Resource Manager `containerregistry` (API Version `2023-06-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-06-01-preview/exportpipelines"
```


### Client Initialization

```go
client := exportpipelines.NewExportPipelinesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ExportPipelinesClient.Create`

```go
ctx := context.TODO()
id := exportpipelines.NewExportPipelineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName", "exportPipelineName")

payload := exportpipelines.ExportPipeline{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ExportPipelinesClient.Delete`

```go
ctx := context.TODO()
id := exportpipelines.NewExportPipelineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName", "exportPipelineName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ExportPipelinesClient.Get`

```go
ctx := context.TODO()
id := exportpipelines.NewExportPipelineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName", "exportPipelineName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ExportPipelinesClient.List`

```go
ctx := context.TODO()
id := exportpipelines.NewRegistryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

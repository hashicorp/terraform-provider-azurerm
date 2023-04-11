
## `github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2021-08-01-preview/importpipelines` Documentation

The `importpipelines` SDK allows for interaction with the Azure Resource Manager Service `containerregistry` (API Version `2021-08-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2021-08-01-preview/importpipelines"
```


### Client Initialization

```go
client := importpipelines.NewImportPipelinesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ImportPipelinesClient.Create`

```go
ctx := context.TODO()
id := importpipelines.NewImportPipelineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue", "importPipelineValue")

payload := importpipelines.ImportPipeline{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ImportPipelinesClient.Delete`

```go
ctx := context.TODO()
id := importpipelines.NewImportPipelineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue", "importPipelineValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ImportPipelinesClient.Get`

```go
ctx := context.TODO()
id := importpipelines.NewImportPipelineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue", "importPipelineValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ImportPipelinesClient.List`

```go
ctx := context.TODO()
id := importpipelines.NewRegistryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

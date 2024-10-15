
## `github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2019-06-01-preview/runs` Documentation

The `runs` SDK allows for interaction with Azure Resource Manager `containerregistry` (API Version `2019-06-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2019-06-01-preview/runs"
```


### Client Initialization

```go
client := runs.NewRunsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `RunsClient.Cancel`

```go
ctx := context.TODO()
id := runs.NewRunID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName", "runId")

if err := client.CancelThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `RunsClient.Get`

```go
ctx := context.TODO()
id := runs.NewRunID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName", "runId")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RunsClient.GetLogSasURL`

```go
ctx := context.TODO()
id := runs.NewRunID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName", "runId")

read, err := client.GetLogSasURL(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RunsClient.List`

```go
ctx := context.TODO()
id := runs.NewRegistryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName")

// alternatively `client.List(ctx, id, runs.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, runs.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `RunsClient.Update`

```go
ctx := context.TODO()
id := runs.NewRunID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName", "runId")

payload := runs.RunUpdateParameters{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


## `github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2019-06-01-preview/taskruns` Documentation

The `taskruns` SDK allows for interaction with Azure Resource Manager `containerregistry` (API Version `2019-06-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2019-06-01-preview/taskruns"
```


### Client Initialization

```go
client := taskruns.NewTaskRunsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `TaskRunsClient.Create`

```go
ctx := context.TODO()
id := taskruns.NewTaskRunID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName", "taskRunName")

payload := taskruns.TaskRun{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `TaskRunsClient.Delete`

```go
ctx := context.TODO()
id := taskruns.NewTaskRunID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName", "taskRunName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `TaskRunsClient.Get`

```go
ctx := context.TODO()
id := taskruns.NewTaskRunID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName", "taskRunName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TaskRunsClient.GetDetails`

```go
ctx := context.TODO()
id := taskruns.NewTaskRunID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName", "taskRunName")

read, err := client.GetDetails(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TaskRunsClient.List`

```go
ctx := context.TODO()
id := taskruns.NewRegistryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `TaskRunsClient.Update`

```go
ctx := context.TODO()
id := taskruns.NewTaskRunID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryName", "taskRunName")

payload := taskruns.TaskRunUpdateParameters{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

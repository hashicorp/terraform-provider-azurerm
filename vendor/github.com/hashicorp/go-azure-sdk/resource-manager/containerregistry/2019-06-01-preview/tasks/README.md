
## `github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2019-06-01-preview/tasks` Documentation

The `tasks` SDK allows for interaction with the Azure Resource Manager Service `containerregistry` (API Version `2019-06-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2019-06-01-preview/tasks"
```


### Client Initialization

```go
client := tasks.NewTasksClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `TasksClient.Create`

```go
ctx := context.TODO()
id := tasks.NewTaskID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue", "taskValue")

payload := tasks.Task{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `TasksClient.Delete`

```go
ctx := context.TODO()
id := tasks.NewTaskID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue", "taskValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `TasksClient.Get`

```go
ctx := context.TODO()
id := tasks.NewTaskID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue", "taskValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TasksClient.GetDetails`

```go
ctx := context.TODO()
id := tasks.NewTaskID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue", "taskValue")

read, err := client.GetDetails(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TasksClient.List`

```go
ctx := context.TODO()
id := tasks.NewRegistryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `TasksClient.Update`

```go
ctx := context.TODO()
id := tasks.NewTaskID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue", "taskValue")

payload := tasks.TaskUpdateParameters{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

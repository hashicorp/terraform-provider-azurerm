
## `github.com/hashicorp/go-azure-sdk/resource-manager/durabletask/2025-11-01/taskhubs` Documentation

The `taskhubs` SDK allows for interaction with Azure Resource Manager `durabletask` (API Version `2025-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/durabletask/2025-11-01/taskhubs"
```


### Client Initialization

```go
client := taskhubs.NewTaskHubsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `TaskHubsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := taskhubs.NewTaskHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "schedulerName", "taskHubName")

payload := taskhubs.TaskHub{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `TaskHubsClient.Delete`

```go
ctx := context.TODO()
id := taskhubs.NewTaskHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "schedulerName", "taskHubName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `TaskHubsClient.Get`

```go
ctx := context.TODO()
id := taskhubs.NewTaskHubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "schedulerName", "taskHubName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TaskHubsClient.ListByScheduler`

```go
ctx := context.TODO()
id := taskhubs.NewSchedulerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "schedulerName")

// alternatively `client.ListByScheduler(ctx, id)` can be used to do batched pagination
items, err := client.ListBySchedulerComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

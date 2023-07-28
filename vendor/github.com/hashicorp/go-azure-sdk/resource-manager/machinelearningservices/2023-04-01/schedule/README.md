
## `github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2023-04-01/schedule` Documentation

The `schedule` SDK allows for interaction with the Azure Resource Manager Service `machinelearningservices` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2023-04-01/schedule"
```


### Client Initialization

```go
client := schedule.NewScheduleClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ScheduleClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := schedule.NewScheduleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "scheduleValue")

payload := schedule.ScheduleResource{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ScheduleClient.Delete`

```go
ctx := context.TODO()
id := schedule.NewScheduleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "scheduleValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ScheduleClient.Get`

```go
ctx := context.TODO()
id := schedule.NewScheduleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "scheduleValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ScheduleClient.List`

```go
ctx := context.TODO()
id := schedule.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue")

// alternatively `client.List(ctx, id, schedule.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, schedule.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

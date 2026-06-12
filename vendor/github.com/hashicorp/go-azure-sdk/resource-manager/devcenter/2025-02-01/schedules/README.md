
## `github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/schedules` Documentation

The `schedules` SDK allows for interaction with Azure Resource Manager `devcenter` (API Version `2025-02-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/schedules"
```


### Client Initialization

```go
client := schedules.NewSchedulesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SchedulesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := schedules.NewScheduleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "projectName", "poolName", "scheduleName")

payload := schedules.Schedule{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `SchedulesClient.Delete`

```go
ctx := context.TODO()
id := schedules.NewScheduleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "projectName", "poolName", "scheduleName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `SchedulesClient.Get`

```go
ctx := context.TODO()
id := schedules.NewScheduleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "projectName", "poolName", "scheduleName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SchedulesClient.ListByPool`

```go
ctx := context.TODO()
id := schedules.NewPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "projectName", "poolName")

// alternatively `client.ListByPool(ctx, id)` can be used to do batched pagination
items, err := client.ListByPoolComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `SchedulesClient.Update`

```go
ctx := context.TODO()
id := schedules.NewScheduleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "projectName", "poolName", "scheduleName")

payload := schedules.ScheduleUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

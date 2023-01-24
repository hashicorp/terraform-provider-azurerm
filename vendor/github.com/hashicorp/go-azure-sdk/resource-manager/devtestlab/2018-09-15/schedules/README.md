
## `github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/schedules` Documentation

The `schedules` SDK allows for interaction with the Azure Resource Manager Service `devtestlab` (API Version `2018-09-15`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/schedules"
```


### Client Initialization

```go
client := schedules.NewSchedulesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SchedulesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := schedules.NewLabScheduleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "labValue", "scheduleValue")

payload := schedules.Schedule{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SchedulesClient.Delete`

```go
ctx := context.TODO()
id := schedules.NewLabScheduleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "labValue", "scheduleValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SchedulesClient.Execute`

```go
ctx := context.TODO()
id := schedules.NewLabScheduleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "labValue", "scheduleValue")

if err := client.ExecuteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `SchedulesClient.Get`

```go
ctx := context.TODO()
id := schedules.NewLabScheduleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "labValue", "scheduleValue")

read, err := client.Get(ctx, id, schedules.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SchedulesClient.List`

```go
ctx := context.TODO()
id := schedules.NewLabID("12345678-1234-9876-4563-123456789012", "example-resource-group", "labValue")

// alternatively `client.List(ctx, id, schedules.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, schedules.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `SchedulesClient.ListApplicable`

```go
ctx := context.TODO()
id := schedules.NewLabScheduleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "labValue", "scheduleValue")

// alternatively `client.ListApplicable(ctx, id)` can be used to do batched pagination
items, err := client.ListApplicableComplete(ctx, id)
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
id := schedules.NewLabScheduleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "labValue", "scheduleValue")

payload := schedules.UpdateResource{
	// ...
}


read, err := client.Update(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

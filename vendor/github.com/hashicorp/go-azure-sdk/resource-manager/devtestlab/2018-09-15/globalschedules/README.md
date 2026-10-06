
## `github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/globalschedules` Documentation

The `globalschedules` SDK allows for interaction with Azure Resource Manager `devtestlab` (API Version `2018-09-15`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/globalschedules"
```


### Client Initialization

```go
client := globalschedules.NewGlobalSchedulesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `GlobalSchedulesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := globalschedules.NewScheduleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "scheduleName")

payload := globalschedules.Schedule{
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


### Example Usage: `GlobalSchedulesClient.Delete`

```go
ctx := context.TODO()
id := globalschedules.NewScheduleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "scheduleName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GlobalSchedulesClient.Execute`

```go
ctx := context.TODO()
id := globalschedules.NewScheduleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "scheduleName")

if err := client.ExecuteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `GlobalSchedulesClient.Get`

```go
ctx := context.TODO()
id := globalschedules.NewScheduleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "scheduleName")

read, err := client.Get(ctx, id, globalschedules.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GlobalSchedulesClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id, globalschedules.DefaultListByResourceGroupOperationOptions())` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id, globalschedules.DefaultListByResourceGroupOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `GlobalSchedulesClient.ListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id, globalschedules.DefaultListBySubscriptionOperationOptions())` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id, globalschedules.DefaultListBySubscriptionOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `GlobalSchedulesClient.Retarget`

```go
ctx := context.TODO()
id := globalschedules.NewScheduleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "scheduleName")

payload := globalschedules.RetargetScheduleProperties{
	// ...
}


if err := client.RetargetThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `GlobalSchedulesClient.Update`

```go
ctx := context.TODO()
id := globalschedules.NewScheduleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "scheduleName")

payload := globalschedules.UpdateResource{
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

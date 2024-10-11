
## `github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2022-06-01-preview/scheduledactions` Documentation

The `scheduledactions` SDK allows for interaction with Azure Resource Manager `costmanagement` (API Version `2022-06-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2022-06-01-preview/scheduledactions"
```


### Client Initialization

```go
client := scheduledactions.NewScheduledActionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ScheduledActionsClient.CheckNameAvailability`

```go
ctx := context.TODO()

payload := scheduledactions.CheckNameAvailabilityRequest{
	// ...
}


read, err := client.CheckNameAvailability(ctx, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ScheduledActionsClient.CheckNameAvailabilityByScope`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

payload := scheduledactions.CheckNameAvailabilityRequest{
	// ...
}


read, err := client.CheckNameAvailabilityByScope(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ScheduledActionsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := scheduledactions.NewScheduledActionID("scheduledActionName")

payload := scheduledactions.ScheduledAction{
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


### Example Usage: `ScheduledActionsClient.CreateOrUpdateByScope`

```go
ctx := context.TODO()
id := scheduledactions.NewScopedScheduledActionID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "scheduledActionName")

payload := scheduledactions.ScheduledAction{
	// ...
}


read, err := client.CreateOrUpdateByScope(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ScheduledActionsClient.Delete`

```go
ctx := context.TODO()
id := scheduledactions.NewScheduledActionID("scheduledActionName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ScheduledActionsClient.DeleteByScope`

```go
ctx := context.TODO()
id := scheduledactions.NewScopedScheduledActionID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "scheduledActionName")

read, err := client.DeleteByScope(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ScheduledActionsClient.Execute`

```go
ctx := context.TODO()
id := scheduledactions.NewScheduledActionID("scheduledActionName")

read, err := client.Execute(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ScheduledActionsClient.ExecuteByScope`

```go
ctx := context.TODO()
id := scheduledactions.NewScopedScheduledActionID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "scheduledActionName")

read, err := client.ExecuteByScope(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ScheduledActionsClient.Get`

```go
ctx := context.TODO()
id := scheduledactions.NewScheduledActionID("scheduledActionName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ScheduledActionsClient.GetByScope`

```go
ctx := context.TODO()
id := scheduledactions.NewScopedScheduledActionID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "scheduledActionName")

read, err := client.GetByScope(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ScheduledActionsClient.List`

```go
ctx := context.TODO()


// alternatively `client.List(ctx, scheduledactions.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, scheduledactions.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ScheduledActionsClient.ListByScope`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

// alternatively `client.ListByScope(ctx, id, scheduledactions.DefaultListByScopeOperationOptions())` can be used to do batched pagination
items, err := client.ListByScopeComplete(ctx, id, scheduledactions.DefaultListByScopeOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


## `github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2025-03-01/scheduledactionoperationgroup` Documentation

The `scheduledactionoperationgroup` SDK allows for interaction with Azure Resource Manager `costmanagement` (API Version `2025-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2025-03-01/scheduledactionoperationgroup"
```


### Client Initialization

```go
client := scheduledactionoperationgroup.NewScheduledActionOperationGroupClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ScheduledActionOperationGroupClient.ScheduledActionsCreateOrUpdateByScope`

```go
ctx := context.TODO()
id := scheduledactionoperationgroup.NewScopedScheduledActionID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "scheduledActionName")

payload := scheduledactionoperationgroup.ScheduledAction{
	// ...
}


read, err := client.ScheduledActionsCreateOrUpdateByScope(ctx, id, payload, scheduledactionoperationgroup.DefaultScheduledActionsCreateOrUpdateByScopeOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ScheduledActionOperationGroupClient.ScheduledActionsDeleteByScope`

```go
ctx := context.TODO()
id := scheduledactionoperationgroup.NewScopedScheduledActionID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "scheduledActionName")

read, err := client.ScheduledActionsDeleteByScope(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ScheduledActionOperationGroupClient.ScheduledActionsGetByScope`

```go
ctx := context.TODO()
id := scheduledactionoperationgroup.NewScopedScheduledActionID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "scheduledActionName")

read, err := client.ScheduledActionsGetByScope(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ScheduledActionOperationGroupClient.ScheduledActionsListByScope`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

// alternatively `client.ScheduledActionsListByScope(ctx, id, scheduledactionoperationgroup.DefaultScheduledActionsListByScopeOperationOptions())` can be used to do batched pagination
items, err := client.ScheduledActionsListByScopeComplete(ctx, id, scheduledactionoperationgroup.DefaultScheduledActionsListByScopeOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ScheduledActionOperationGroupClient.ScheduledActionsRunByScope`

```go
ctx := context.TODO()
id := scheduledactionoperationgroup.NewScopedScheduledActionID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "scheduledActionName")

read, err := client.ScheduledActionsRunByScope(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

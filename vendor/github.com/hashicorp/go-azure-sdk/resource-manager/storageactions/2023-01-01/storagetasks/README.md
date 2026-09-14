
## `github.com/hashicorp/go-azure-sdk/resource-manager/storageactions/2023-01-01/storagetasks` Documentation

The `storagetasks` SDK allows for interaction with Azure Resource Manager `storageactions` (API Version `2023-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/storageactions/2023-01-01/storagetasks"
```


### Client Initialization

```go
client := storagetasks.NewStorageTasksClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `StorageTasksClient.Create`

```go
ctx := context.TODO()
id := storagetasks.NewStorageTaskID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageTaskName")

payload := storagetasks.StorageTask{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `StorageTasksClient.Delete`

```go
ctx := context.TODO()
id := storagetasks.NewStorageTaskID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageTaskName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `StorageTasksClient.Get`

```go
ctx := context.TODO()
id := storagetasks.NewStorageTaskID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageTaskName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StorageTasksClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `StorageTasksClient.ListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `StorageTasksClient.ReportList`

```go
ctx := context.TODO()
id := storagetasks.NewStorageTaskID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageTaskName")

// alternatively `client.ReportList(ctx, id, storagetasks.DefaultReportListOperationOptions())` can be used to do batched pagination
items, err := client.ReportListComplete(ctx, id, storagetasks.DefaultReportListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `StorageTasksClient.StorageTaskAssignmentList`

```go
ctx := context.TODO()
id := storagetasks.NewStorageTaskID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageTaskName")

// alternatively `client.StorageTaskAssignmentList(ctx, id, storagetasks.DefaultStorageTaskAssignmentListOperationOptions())` can be used to do batched pagination
items, err := client.StorageTaskAssignmentListComplete(ctx, id, storagetasks.DefaultStorageTaskAssignmentListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `StorageTasksClient.Update`

```go
ctx := context.TODO()
id := storagetasks.NewStorageTaskID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageTaskName")

payload := storagetasks.StorageTaskUpdateParameters{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

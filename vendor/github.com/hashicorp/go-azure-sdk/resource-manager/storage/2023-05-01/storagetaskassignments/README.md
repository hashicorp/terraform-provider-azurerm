
## `github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-05-01/storagetaskassignments` Documentation

The `storagetaskassignments` SDK allows for interaction with Azure Resource Manager `storage` (API Version `2023-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-05-01/storagetaskassignments"
```


### Client Initialization

```go
client := storagetaskassignments.NewStorageTaskAssignmentsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `StorageTaskAssignmentsClient.Create`

```go
ctx := context.TODO()
id := storagetaskassignments.NewStorageTaskAssignmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName", "storageTaskAssignmentName")

payload := storagetaskassignments.StorageTaskAssignment{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `StorageTaskAssignmentsClient.Delete`

```go
ctx := context.TODO()
id := storagetaskassignments.NewStorageTaskAssignmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName", "storageTaskAssignmentName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `StorageTaskAssignmentsClient.Get`

```go
ctx := context.TODO()
id := storagetaskassignments.NewStorageTaskAssignmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName", "storageTaskAssignmentName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StorageTaskAssignmentsClient.InstancesReportList`

```go
ctx := context.TODO()
id := commonids.NewStorageAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName")

// alternatively `client.InstancesReportList(ctx, id, storagetaskassignments.DefaultInstancesReportListOperationOptions())` can be used to do batched pagination
items, err := client.InstancesReportListComplete(ctx, id, storagetaskassignments.DefaultInstancesReportListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `StorageTaskAssignmentsClient.List`

```go
ctx := context.TODO()
id := commonids.NewStorageAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName")

// alternatively `client.List(ctx, id, storagetaskassignments.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, storagetaskassignments.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `StorageTaskAssignmentsClient.StorageTaskAssignmentInstancesReportList`

```go
ctx := context.TODO()
id := storagetaskassignments.NewStorageTaskAssignmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName", "storageTaskAssignmentName")

// alternatively `client.StorageTaskAssignmentInstancesReportList(ctx, id, storagetaskassignments.DefaultStorageTaskAssignmentInstancesReportListOperationOptions())` can be used to do batched pagination
items, err := client.StorageTaskAssignmentInstancesReportListComplete(ctx, id, storagetaskassignments.DefaultStorageTaskAssignmentInstancesReportListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `StorageTaskAssignmentsClient.Update`

```go
ctx := context.TODO()
id := storagetaskassignments.NewStorageTaskAssignmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName", "storageTaskAssignmentName")

payload := storagetaskassignments.StorageTaskAssignmentUpdateParameters{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

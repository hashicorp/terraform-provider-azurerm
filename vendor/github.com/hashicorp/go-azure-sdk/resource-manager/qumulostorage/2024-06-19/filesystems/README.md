
## `github.com/hashicorp/go-azure-sdk/resource-manager/qumulostorage/2024-06-19/filesystems` Documentation

The `filesystems` SDK allows for interaction with Azure Resource Manager `qumulostorage` (API Version `2024-06-19`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/qumulostorage/2024-06-19/filesystems"
```


### Client Initialization

```go
client := filesystems.NewFileSystemsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `FileSystemsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := filesystems.NewFileSystemID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fileSystemName")

payload := filesystems.LiftrBaseStorageFileSystemResource{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `FileSystemsClient.Delete`

```go
ctx := context.TODO()
id := filesystems.NewFileSystemID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fileSystemName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `FileSystemsClient.Get`

```go
ctx := context.TODO()
id := filesystems.NewFileSystemID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fileSystemName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FileSystemsClient.ListByResourceGroup`

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


### Example Usage: `FileSystemsClient.ListBySubscription`

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


### Example Usage: `FileSystemsClient.Update`

```go
ctx := context.TODO()
id := filesystems.NewFileSystemID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fileSystemName")

payload := filesystems.LiftrBaseStorageFileSystemResourceUpdate{
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

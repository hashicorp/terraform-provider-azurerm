
## `github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2023-05-01/storagetargets` Documentation

The `storagetargets` SDK allows for interaction with the Azure Resource Manager Service `storagecache` (API Version `2023-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2023-05-01/storagetargets"
```


### Client Initialization

```go
client := storagetargets.NewStorageTargetsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `StorageTargetsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := storagetargets.NewStorageTargetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cacheValue", "storageTargetValue")

payload := storagetargets.StorageTarget{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `StorageTargetsClient.Delete`

```go
ctx := context.TODO()
id := storagetargets.NewStorageTargetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cacheValue", "storageTargetValue")

if err := client.DeleteThenPoll(ctx, id, storagetargets.DefaultDeleteOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `StorageTargetsClient.DnsRefresh`

```go
ctx := context.TODO()
id := storagetargets.NewStorageTargetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cacheValue", "storageTargetValue")

if err := client.DnsRefreshThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `StorageTargetsClient.Get`

```go
ctx := context.TODO()
id := storagetargets.NewStorageTargetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cacheValue", "storageTargetValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StorageTargetsClient.ListByCache`

```go
ctx := context.TODO()
id := storagetargets.NewCacheID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cacheValue")

// alternatively `client.ListByCache(ctx, id)` can be used to do batched pagination
items, err := client.ListByCacheComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `StorageTargetsClient.RestoreDefaults`

```go
ctx := context.TODO()
id := storagetargets.NewStorageTargetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cacheValue", "storageTargetValue")

if err := client.RestoreDefaultsThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `StorageTargetsClient.StorageTargetFlush`

```go
ctx := context.TODO()
id := storagetargets.NewStorageTargetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cacheValue", "storageTargetValue")

if err := client.StorageTargetFlushThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `StorageTargetsClient.StorageTargetInvalidate`

```go
ctx := context.TODO()
id := storagetargets.NewStorageTargetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cacheValue", "storageTargetValue")

if err := client.StorageTargetInvalidateThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `StorageTargetsClient.StorageTargetResume`

```go
ctx := context.TODO()
id := storagetargets.NewStorageTargetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cacheValue", "storageTargetValue")

if err := client.StorageTargetResumeThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `StorageTargetsClient.StorageTargetSuspend`

```go
ctx := context.TODO()
id := storagetargets.NewStorageTargetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cacheValue", "storageTargetValue")

if err := client.StorageTargetSuspendThenPoll(ctx, id); err != nil {
	// handle the error
}
```

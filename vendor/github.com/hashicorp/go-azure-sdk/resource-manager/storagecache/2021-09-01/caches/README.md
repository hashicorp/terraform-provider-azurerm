
## `github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2021-09-01/caches` Documentation

The `caches` SDK allows for interaction with the Azure Resource Manager Service `storagecache` (API Version `2021-09-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2021-09-01/caches"
```


### Client Initialization

```go
client := caches.NewCachesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `CachesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := caches.NewCacheID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cacheValue")

payload := caches.Cache{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CachesClient.DebugInfo`

```go
ctx := context.TODO()
id := caches.NewCacheID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cacheValue")

if err := client.DebugInfoThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CachesClient.Delete`

```go
ctx := context.TODO()
id := caches.NewCacheID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cacheValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CachesClient.Flush`

```go
ctx := context.TODO()
id := caches.NewCacheID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cacheValue")

if err := client.FlushThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CachesClient.Get`

```go
ctx := context.TODO()
id := caches.NewCacheID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cacheValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CachesClient.List`

```go
ctx := context.TODO()
id := caches.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `CachesClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := caches.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `CachesClient.Start`

```go
ctx := context.TODO()
id := caches.NewCacheID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cacheValue")

if err := client.StartThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CachesClient.Stop`

```go
ctx := context.TODO()
id := caches.NewCacheID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cacheValue")

if err := client.StopThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CachesClient.Update`

```go
ctx := context.TODO()
id := caches.NewCacheID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cacheValue")

payload := caches.Cache{
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


### Example Usage: `CachesClient.UpgradeFirmware`

```go
ctx := context.TODO()
id := caches.NewCacheID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cacheValue")

if err := client.UpgradeFirmwareThenPoll(ctx, id); err != nil {
	// handle the error
}
```

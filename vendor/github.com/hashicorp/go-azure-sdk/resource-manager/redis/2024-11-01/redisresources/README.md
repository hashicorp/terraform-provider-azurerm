
## `github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-11-01/redisresources` Documentation

The `redisresources` SDK allows for interaction with Azure Resource Manager `redis` (API Version `2024-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-11-01/redisresources"
```


### Client Initialization

```go
client := redisresources.NewRedisResourcesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `RedisResourcesClient.PrivateLinkResourcesListByRedisCache`

```go
ctx := context.TODO()
id := redisresources.NewRediID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisName")

// alternatively `client.PrivateLinkResourcesListByRedisCache(ctx, id)` can be used to do batched pagination
items, err := client.PrivateLinkResourcesListByRedisCacheComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `RedisResourcesClient.RedisCreate`

```go
ctx := context.TODO()
id := redisresources.NewRediID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisName")

payload := redisresources.RedisCreateParameters{
	// ...
}


if err := client.RedisCreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `RedisResourcesClient.RedisDelete`

```go
ctx := context.TODO()
id := redisresources.NewRediID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisName")

if err := client.RedisDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `RedisResourcesClient.RedisExportData`

```go
ctx := context.TODO()
id := redisresources.NewRediID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisName")

payload := redisresources.ExportRDBParameters{
	// ...
}


if err := client.RedisExportDataThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `RedisResourcesClient.RedisFlushCache`

```go
ctx := context.TODO()
id := redisresources.NewRediID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisName")

if err := client.RedisFlushCacheThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `RedisResourcesClient.RedisForceReboot`

```go
ctx := context.TODO()
id := redisresources.NewRediID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisName")

payload := redisresources.RedisRebootParameters{
	// ...
}


read, err := client.RedisForceReboot(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RedisResourcesClient.RedisGet`

```go
ctx := context.TODO()
id := redisresources.NewRediID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisName")

read, err := client.RedisGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RedisResourcesClient.RedisImportData`

```go
ctx := context.TODO()
id := redisresources.NewRediID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisName")

payload := redisresources.ImportRDBParameters{
	// ...
}


if err := client.RedisImportDataThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `RedisResourcesClient.RedisListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.RedisListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.RedisListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `RedisResourcesClient.RedisListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.RedisListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.RedisListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `RedisResourcesClient.RedisListKeys`

```go
ctx := context.TODO()
id := redisresources.NewRediID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisName")

read, err := client.RedisListKeys(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RedisResourcesClient.RedisListUpgradeNotifications`

```go
ctx := context.TODO()
id := redisresources.NewRediID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisName")

// alternatively `client.RedisListUpgradeNotifications(ctx, id, redisresources.DefaultRedisListUpgradeNotificationsOperationOptions())` can be used to do batched pagination
items, err := client.RedisListUpgradeNotificationsComplete(ctx, id, redisresources.DefaultRedisListUpgradeNotificationsOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `RedisResourcesClient.RedisRegenerateKey`

```go
ctx := context.TODO()
id := redisresources.NewRediID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisName")

payload := redisresources.RedisRegenerateKeyParameters{
	// ...
}


read, err := client.RedisRegenerateKey(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RedisResourcesClient.RedisUpdate`

```go
ctx := context.TODO()
id := redisresources.NewRediID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisName")

payload := redisresources.RedisUpdateParameters{
	// ...
}


if err := client.RedisUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


## `github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-11-01/redispatchschedules` Documentation

The `redispatchschedules` SDK allows for interaction with Azure Resource Manager `redis` (API Version `2024-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-11-01/redispatchschedules"
```


### Client Initialization

```go
client := redispatchschedules.NewRedisPatchSchedulesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `RedisPatchSchedulesClient.PatchSchedulesCreateOrUpdate`

```go
ctx := context.TODO()
id := redispatchschedules.NewRediID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisName")

payload := redispatchschedules.RedisPatchSchedule{
	// ...
}


read, err := client.PatchSchedulesCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RedisPatchSchedulesClient.PatchSchedulesDelete`

```go
ctx := context.TODO()
id := redispatchschedules.NewRediID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisName")

read, err := client.PatchSchedulesDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RedisPatchSchedulesClient.PatchSchedulesGet`

```go
ctx := context.TODO()
id := redispatchschedules.NewRediID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisName")

read, err := client.PatchSchedulesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RedisPatchSchedulesClient.PatchSchedulesListByRedisResource`

```go
ctx := context.TODO()
id := redispatchschedules.NewRediID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisName")

// alternatively `client.PatchSchedulesListByRedisResource(ctx, id)` can be used to do batched pagination
items, err := client.PatchSchedulesListByRedisResourceComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

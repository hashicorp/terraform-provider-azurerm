
## `github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-03-01/patchschedules` Documentation

The `patchschedules` SDK allows for interaction with Azure Resource Manager `redis` (API Version `2024-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-03-01/patchschedules"
```


### Client Initialization

```go
client := patchschedules.NewPatchSchedulesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PatchSchedulesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := patchschedules.NewRediID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cacheName")

payload := patchschedules.RedisPatchSchedule{
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


### Example Usage: `PatchSchedulesClient.Delete`

```go
ctx := context.TODO()
id := patchschedules.NewRediID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cacheName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PatchSchedulesClient.Get`

```go
ctx := context.TODO()
id := patchschedules.NewRediID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cacheName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PatchSchedulesClient.ListByRedisResource`

```go
ctx := context.TODO()
id := patchschedules.NewRediID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cacheName")

// alternatively `client.ListByRedisResource(ctx, id)` can be used to do batched pagination
items, err := client.ListByRedisResourceComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

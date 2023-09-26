
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/cache` Documentation

The `cache` SDK allows for interaction with the Azure Resource Manager Service `apimanagement` (API Version `2021-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/cache"
```


### Client Initialization

```go
client := cache.NewCacheClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `CacheClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := cache.NewCacheID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "cacheIdValue")

payload := cache.CacheContract{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, cache.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CacheClient.Delete`

```go
ctx := context.TODO()
id := cache.NewCacheID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "cacheIdValue")

read, err := client.Delete(ctx, id, cache.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CacheClient.Get`

```go
ctx := context.TODO()
id := cache.NewCacheID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "cacheIdValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CacheClient.GetEntityTag`

```go
ctx := context.TODO()
id := cache.NewCacheID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "cacheIdValue")

read, err := client.GetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CacheClient.ListByService`

```go
ctx := context.TODO()
id := cache.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue")

// alternatively `client.ListByService(ctx, id, cache.DefaultListByServiceOperationOptions())` can be used to do batched pagination
items, err := client.ListByServiceComplete(ctx, id, cache.DefaultListByServiceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `CacheClient.Update`

```go
ctx := context.TODO()
id := cache.NewCacheID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "cacheIdValue")

payload := cache.CacheUpdateParameters{
	// ...
}


read, err := client.Update(ctx, id, payload, cache.DefaultUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

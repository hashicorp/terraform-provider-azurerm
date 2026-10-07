
## `github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-11-01/redis` Documentation

The `redis` SDK allows for interaction with Azure Resource Manager `redis` (API Version `2024-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-11-01/redis"
```


### Client Initialization

```go
client := redis.NewRedisClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `RedisClient.AsyncOperationStatusGet`

```go
ctx := context.TODO()
id := redis.NewAsyncOperationID("12345678-1234-9876-4563-123456789012", "locationName", "operationId")

read, err := client.AsyncOperationStatusGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RedisClient.RedisCheckNameAvailability`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

payload := redis.CheckNameAvailabilityParameters{
	// ...
}


read, err := client.RedisCheckNameAvailability(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

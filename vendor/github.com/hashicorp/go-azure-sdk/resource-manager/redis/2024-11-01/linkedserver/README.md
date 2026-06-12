
## `github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-11-01/linkedserver` Documentation

The `linkedserver` SDK allows for interaction with Azure Resource Manager `redis` (API Version `2024-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-11-01/linkedserver"
```


### Client Initialization

```go
client := linkedserver.NewLinkedServerClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `LinkedServerClient.Create`

```go
ctx := context.TODO()
id := linkedserver.NewLinkedServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisName", "linkedServerName")

payload := linkedserver.RedisLinkedServerCreateParameters{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `LinkedServerClient.Delete`

```go
ctx := context.TODO()
id := linkedserver.NewLinkedServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisName", "linkedServerName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `LinkedServerClient.Get`

```go
ctx := context.TODO()
id := linkedserver.NewLinkedServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisName", "linkedServerName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LinkedServerClient.List`

```go
ctx := context.TODO()
id := linkedserver.NewRediID("12345678-1234-9876-4563-123456789012", "example-resource-group", "redisName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

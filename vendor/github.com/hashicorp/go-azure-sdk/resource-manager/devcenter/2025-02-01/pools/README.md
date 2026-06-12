
## `github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/pools` Documentation

The `pools` SDK allows for interaction with Azure Resource Manager `devcenter` (API Version `2025-02-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/pools"
```


### Client Initialization

```go
client := pools.NewPoolsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PoolsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := pools.NewPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "projectName", "poolName")

payload := pools.Pool{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `PoolsClient.Delete`

```go
ctx := context.TODO()
id := pools.NewPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "projectName", "poolName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `PoolsClient.Get`

```go
ctx := context.TODO()
id := pools.NewPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "projectName", "poolName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PoolsClient.ListByProject`

```go
ctx := context.TODO()
id := pools.NewProjectID("12345678-1234-9876-4563-123456789012", "example-resource-group", "projectName")

// alternatively `client.ListByProject(ctx, id)` can be used to do batched pagination
items, err := client.ListByProjectComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PoolsClient.RunHealthChecks`

```go
ctx := context.TODO()
id := pools.NewPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "projectName", "poolName")

if err := client.RunHealthChecksThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `PoolsClient.Update`

```go
ctx := context.TODO()
id := pools.NewPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "projectName", "poolName")

payload := pools.PoolUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


## `github.com/hashicorp/go-azure-sdk/resource-manager/horizondb/2026-10-01-preview/pools` Documentation

The `pools` SDK allows for interaction with Azure Resource Manager `horizondb` (API Version `2026-10-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/horizondb/2026-10-01-preview/pools"
```


### Client Initialization

```go
client := pools.NewPoolsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PoolsClient.CreateOrReplace`

```go
ctx := context.TODO()
id := pools.NewPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "poolName")

payload := pools.Pool{
	// ...
}


if err := client.CreateOrReplaceThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `PoolsClient.Delete`

```go
ctx := context.TODO()
id := pools.NewPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "poolName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `PoolsClient.Get`

```go
ctx := context.TODO()
id := pools.NewPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "poolName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PoolsClient.List`

```go
ctx := context.TODO()
id := pools.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PoolsClient.Restart`

```go
ctx := context.TODO()
id := pools.NewPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "poolName")

if err := client.RestartThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `PoolsClient.Start`

```go
ctx := context.TODO()
id := pools.NewPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "poolName")

if err := client.StartThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `PoolsClient.Stop`

```go
ctx := context.TODO()
id := pools.NewPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "poolName")

if err := client.StopThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `PoolsClient.Update`

```go
ctx := context.TODO()
id := pools.NewPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterName", "poolName")

payload := pools.PoolForPatchUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

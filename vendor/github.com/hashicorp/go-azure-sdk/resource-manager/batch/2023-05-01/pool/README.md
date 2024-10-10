
## `github.com/hashicorp/go-azure-sdk/resource-manager/batch/2023-05-01/pool` Documentation

The `pool` SDK allows for interaction with Azure Resource Manager `batch` (API Version `2023-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/batch/2023-05-01/pool"
```


### Client Initialization

```go
client := pool.NewPoolClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PoolClient.Create`

```go
ctx := context.TODO()
id := pool.NewPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "batchAccountName", "poolName")

payload := pool.Pool{
	// ...
}


read, err := client.Create(ctx, id, payload, pool.DefaultCreateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PoolClient.Delete`

```go
ctx := context.TODO()
id := pool.NewPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "batchAccountName", "poolName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `PoolClient.DisableAutoScale`

```go
ctx := context.TODO()
id := pool.NewPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "batchAccountName", "poolName")

read, err := client.DisableAutoScale(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PoolClient.Get`

```go
ctx := context.TODO()
id := pool.NewPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "batchAccountName", "poolName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PoolClient.ListByBatchAccount`

```go
ctx := context.TODO()
id := pool.NewBatchAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "batchAccountName")

// alternatively `client.ListByBatchAccount(ctx, id, pool.DefaultListByBatchAccountOperationOptions())` can be used to do batched pagination
items, err := client.ListByBatchAccountComplete(ctx, id, pool.DefaultListByBatchAccountOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PoolClient.StopResize`

```go
ctx := context.TODO()
id := pool.NewPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "batchAccountName", "poolName")

read, err := client.StopResize(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PoolClient.Update`

```go
ctx := context.TODO()
id := pool.NewPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "batchAccountName", "poolName")

payload := pool.Pool{
	// ...
}


read, err := client.Update(ctx, id, payload, pool.DefaultUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

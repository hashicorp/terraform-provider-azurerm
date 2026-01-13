
## `github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/pools` Documentation

The `pools` SDK allows for interaction with Azure Resource Manager `batch` (API Version `2024-07-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/pools"
```


### Client Initialization

```go
client := pools.NewPoolsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PoolsClient.PoolCreate`

```go
ctx := context.TODO()
id := pools.NewPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "batchAccountName", "poolName")

payload := pools.Pool{
	// ...
}


read, err := client.PoolCreate(ctx, id, payload, pools.DefaultPoolCreateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PoolsClient.PoolDelete`

```go
ctx := context.TODO()
id := pools.NewPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "batchAccountName", "poolName")

if err := client.PoolDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `PoolsClient.PoolDisableAutoScale`

```go
ctx := context.TODO()
id := pools.NewPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "batchAccountName", "poolName")

read, err := client.PoolDisableAutoScale(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PoolsClient.PoolGet`

```go
ctx := context.TODO()
id := pools.NewPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "batchAccountName", "poolName")

read, err := client.PoolGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PoolsClient.PoolListByBatchAccount`

```go
ctx := context.TODO()
id := pools.NewBatchAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "batchAccountName")

// alternatively `client.PoolListByBatchAccount(ctx, id, pools.DefaultPoolListByBatchAccountOperationOptions())` can be used to do batched pagination
items, err := client.PoolListByBatchAccountComplete(ctx, id, pools.DefaultPoolListByBatchAccountOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PoolsClient.PoolStopResize`

```go
ctx := context.TODO()
id := pools.NewPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "batchAccountName", "poolName")

read, err := client.PoolStopResize(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PoolsClient.PoolUpdate`

```go
ctx := context.TODO()
id := pools.NewPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "batchAccountName", "poolName")

payload := pools.Pool{
	// ...
}


read, err := client.PoolUpdate(ctx, id, payload, pools.DefaultPoolUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

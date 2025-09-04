
## `github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-01-01/capacitypools` Documentation

The `capacitypools` SDK allows for interaction with Azure Resource Manager `netapp` (API Version `2025-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-01-01/capacitypools"
```


### Client Initialization

```go
client := capacitypools.NewCapacityPoolsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `CapacityPoolsClient.PoolsCreateOrUpdate`

```go
ctx := context.TODO()
id := capacitypools.NewCapacityPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "capacityPoolName")

payload := capacitypools.CapacityPool{
	// ...
}


if err := client.PoolsCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CapacityPoolsClient.PoolsDelete`

```go
ctx := context.TODO()
id := capacitypools.NewCapacityPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "capacityPoolName")

if err := client.PoolsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CapacityPoolsClient.PoolsGet`

```go
ctx := context.TODO()
id := capacitypools.NewCapacityPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "capacityPoolName")

read, err := client.PoolsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CapacityPoolsClient.PoolsList`

```go
ctx := context.TODO()
id := capacitypools.NewNetAppAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName")

// alternatively `client.PoolsList(ctx, id)` can be used to do batched pagination
items, err := client.PoolsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `CapacityPoolsClient.PoolsUpdate`

```go
ctx := context.TODO()
id := capacitypools.NewCapacityPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "capacityPoolName")

payload := capacitypools.CapacityPoolPatch{
	// ...
}


if err := client.PoolsUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

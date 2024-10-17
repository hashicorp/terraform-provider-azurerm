
## `github.com/hashicorp/go-azure-sdk/resource-manager/devopsinfrastructure/2024-04-04-preview/pools` Documentation

The `pools` SDK allows for interaction with Azure Resource Manager `devopsinfrastructure` (API Version `2024-04-04-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/devopsinfrastructure/2024-04-04-preview/pools"
```


### Client Initialization

```go
client := pools.NewPoolsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PoolsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := pools.NewPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "poolName")

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
id := pools.NewPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "poolName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `PoolsClient.Get`

```go
ctx := context.TODO()
id := pools.NewPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "poolName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PoolsClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PoolsClient.ListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PoolsClient.Update`

```go
ctx := context.TODO()
id := pools.NewPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "poolName")

payload := pools.PoolUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

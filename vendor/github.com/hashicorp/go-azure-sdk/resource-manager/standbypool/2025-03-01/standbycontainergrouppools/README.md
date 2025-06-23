
## `github.com/hashicorp/go-azure-sdk/resource-manager/standbypool/2025-03-01/standbycontainergrouppools` Documentation

The `standbycontainergrouppools` SDK allows for interaction with Azure Resource Manager `standbypool` (API Version `2025-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/standbypool/2025-03-01/standbycontainergrouppools"
```


### Client Initialization

```go
client := standbycontainergrouppools.NewStandbyContainerGroupPoolsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `StandbyContainerGroupPoolsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := standbycontainergrouppools.NewStandbyContainerGroupPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "standbyContainerGroupPoolName")

payload := standbycontainergrouppools.StandbyContainerGroupPoolResource{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `StandbyContainerGroupPoolsClient.Delete`

```go
ctx := context.TODO()
id := standbycontainergrouppools.NewStandbyContainerGroupPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "standbyContainerGroupPoolName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `StandbyContainerGroupPoolsClient.Get`

```go
ctx := context.TODO()
id := standbycontainergrouppools.NewStandbyContainerGroupPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "standbyContainerGroupPoolName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StandbyContainerGroupPoolsClient.ListByResourceGroup`

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


### Example Usage: `StandbyContainerGroupPoolsClient.ListBySubscription`

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


### Example Usage: `StandbyContainerGroupPoolsClient.Update`

```go
ctx := context.TODO()
id := standbycontainergrouppools.NewStandbyContainerGroupPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "standbyContainerGroupPoolName")

payload := standbycontainergrouppools.StandbyContainerGroupPoolResourceUpdate{
	// ...
}


read, err := client.Update(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

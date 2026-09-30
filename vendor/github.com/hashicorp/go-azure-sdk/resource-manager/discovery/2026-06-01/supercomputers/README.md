
## `github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/supercomputers` Documentation

The `supercomputers` SDK allows for interaction with Azure Resource Manager `discovery` (API Version `2026-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/supercomputers"
```


### Client Initialization

```go
client := supercomputers.NewSupercomputersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SupercomputersClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := supercomputers.NewSupercomputerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "supercomputerName")

payload := supercomputers.Supercomputer{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `SupercomputersClient.Delete`

```go
ctx := context.TODO()
id := supercomputers.NewSupercomputerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "supercomputerName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `SupercomputersClient.Get`

```go
ctx := context.TODO()
id := supercomputers.NewSupercomputerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "supercomputerName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SupercomputersClient.ListByResourceGroup`

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


### Example Usage: `SupercomputersClient.ListBySubscription`

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


### Example Usage: `SupercomputersClient.Update`

```go
ctx := context.TODO()
id := supercomputers.NewSupercomputerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "supercomputerName")

payload := supercomputers.SupercomputerUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

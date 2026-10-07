
## `github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/networkanchors` Documentation

The `networkanchors` SDK allows for interaction with Azure Resource Manager `oracledatabase` (API Version `2025-09-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/networkanchors"
```


### Client Initialization

```go
client := networkanchors.NewNetworkAnchorsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `NetworkAnchorsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := networkanchors.NewNetworkAnchorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkAnchorName")

payload := networkanchors.NetworkAnchor{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `NetworkAnchorsClient.Delete`

```go
ctx := context.TODO()
id := networkanchors.NewNetworkAnchorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkAnchorName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `NetworkAnchorsClient.Get`

```go
ctx := context.TODO()
id := networkanchors.NewNetworkAnchorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkAnchorName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NetworkAnchorsClient.ListByResourceGroup`

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


### Example Usage: `NetworkAnchorsClient.ListBySubscription`

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


### Example Usage: `NetworkAnchorsClient.Update`

```go
ctx := context.TODO()
id := networkanchors.NewNetworkAnchorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkAnchorName")

payload := networkanchors.NetworkAnchorUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

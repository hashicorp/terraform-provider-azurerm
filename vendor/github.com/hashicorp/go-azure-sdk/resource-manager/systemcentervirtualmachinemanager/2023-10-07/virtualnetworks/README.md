
## `github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/virtualnetworks` Documentation

The `virtualnetworks` SDK allows for interaction with Azure Resource Manager `systemcentervirtualmachinemanager` (API Version `2023-10-07`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/virtualnetworks"
```


### Client Initialization

```go
client := virtualnetworks.NewVirtualNetworksClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `VirtualNetworksClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := virtualnetworks.NewVirtualNetworkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualNetworkName")

payload := virtualnetworks.VirtualNetwork{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualNetworksClient.Delete`

```go
ctx := context.TODO()
id := virtualnetworks.NewVirtualNetworkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualNetworkName")

if err := client.DeleteThenPoll(ctx, id, virtualnetworks.DefaultDeleteOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualNetworksClient.Get`

```go
ctx := context.TODO()
id := virtualnetworks.NewVirtualNetworkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualNetworkName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualNetworksClient.ListByResourceGroup`

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


### Example Usage: `VirtualNetworksClient.ListBySubscription`

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


### Example Usage: `VirtualNetworksClient.Update`

```go
ctx := context.TODO()
id := virtualnetworks.NewVirtualNetworkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualNetworkName")

payload := virtualnetworks.VirtualNetworkTagsUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

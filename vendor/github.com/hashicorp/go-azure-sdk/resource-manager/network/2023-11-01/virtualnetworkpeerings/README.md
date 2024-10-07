
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/virtualnetworkpeerings` Documentation

The `virtualnetworkpeerings` SDK allows for interaction with Azure Resource Manager `network` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/virtualnetworkpeerings"
```


### Client Initialization

```go
client := virtualnetworkpeerings.NewVirtualNetworkPeeringsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `VirtualNetworkPeeringsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := virtualnetworkpeerings.NewVirtualNetworkPeeringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualNetworkName", "virtualNetworkPeeringName")

payload := virtualnetworkpeerings.VirtualNetworkPeering{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload, virtualnetworkpeerings.DefaultCreateOrUpdateOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualNetworkPeeringsClient.Delete`

```go
ctx := context.TODO()
id := virtualnetworkpeerings.NewVirtualNetworkPeeringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualNetworkName", "virtualNetworkPeeringName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualNetworkPeeringsClient.Get`

```go
ctx := context.TODO()
id := virtualnetworkpeerings.NewVirtualNetworkPeeringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualNetworkName", "virtualNetworkPeeringName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualNetworkPeeringsClient.List`

```go
ctx := context.TODO()
id := commonids.NewVirtualNetworkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualNetworkName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/virtualrouterpeerings` Documentation

The `virtualrouterpeerings` SDK allows for interaction with the Azure Resource Manager Service `network` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/virtualrouterpeerings"
```


### Client Initialization

```go
client := virtualrouterpeerings.NewVirtualRouterPeeringsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `VirtualRouterPeeringsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := virtualrouterpeerings.NewVirtualRouterPeeringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualRouterValue", "peeringValue")

payload := virtualrouterpeerings.VirtualRouterPeering{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualRouterPeeringsClient.Delete`

```go
ctx := context.TODO()
id := virtualrouterpeerings.NewVirtualRouterPeeringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualRouterValue", "peeringValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualRouterPeeringsClient.Get`

```go
ctx := context.TODO()
id := virtualrouterpeerings.NewVirtualRouterPeeringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualRouterValue", "peeringValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualRouterPeeringsClient.List`

```go
ctx := context.TODO()
id := virtualrouterpeerings.NewVirtualRouterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualRouterValue")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

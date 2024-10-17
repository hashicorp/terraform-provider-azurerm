
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/peerexpressroutecircuitconnections` Documentation

The `peerexpressroutecircuitconnections` SDK allows for interaction with Azure Resource Manager `network` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/peerexpressroutecircuitconnections"
```


### Client Initialization

```go
client := peerexpressroutecircuitconnections.NewPeerExpressRouteCircuitConnectionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PeerExpressRouteCircuitConnectionsClient.Get`

```go
ctx := context.TODO()
id := peerexpressroutecircuitconnections.NewPeerConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRouteCircuitName", "peeringName", "peerConnectionName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PeerExpressRouteCircuitConnectionsClient.List`

```go
ctx := context.TODO()
id := commonids.NewExpressRouteCircuitPeeringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRouteCircuitName", "peeringName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

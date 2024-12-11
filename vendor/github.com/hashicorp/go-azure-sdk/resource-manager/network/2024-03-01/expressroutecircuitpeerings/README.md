
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-03-01/expressroutecircuitpeerings` Documentation

The `expressroutecircuitpeerings` SDK allows for interaction with Azure Resource Manager `network` (API Version `2024-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-03-01/expressroutecircuitpeerings"
```


### Client Initialization

```go
client := expressroutecircuitpeerings.NewExpressRouteCircuitPeeringsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ExpressRouteCircuitPeeringsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := commonids.NewExpressRouteCircuitPeeringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRouteCircuitName", "peeringName")

payload := expressroutecircuitpeerings.ExpressRouteCircuitPeering{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ExpressRouteCircuitPeeringsClient.Delete`

```go
ctx := context.TODO()
id := commonids.NewExpressRouteCircuitPeeringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRouteCircuitName", "peeringName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ExpressRouteCircuitPeeringsClient.Get`

```go
ctx := context.TODO()
id := commonids.NewExpressRouteCircuitPeeringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRouteCircuitName", "peeringName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ExpressRouteCircuitPeeringsClient.List`

```go
ctx := context.TODO()
id := expressroutecircuitpeerings.NewExpressRouteCircuitID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRouteCircuitName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

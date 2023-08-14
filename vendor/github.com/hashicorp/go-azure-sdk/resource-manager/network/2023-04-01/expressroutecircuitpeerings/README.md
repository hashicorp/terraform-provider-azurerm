
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/expressroutecircuitpeerings` Documentation

The `expressroutecircuitpeerings` SDK allows for interaction with the Azure Resource Manager Service `network` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/expressroutecircuitpeerings"
```


### Client Initialization

```go
client := expressroutecircuitpeerings.NewExpressRouteCircuitPeeringsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ExpressRouteCircuitPeeringsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := expressroutecircuitpeerings.NewExpressRouteCircuitPeeringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRouteCircuitValue", "peeringValue")

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
id := expressroutecircuitpeerings.NewExpressRouteCircuitPeeringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRouteCircuitValue", "peeringValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ExpressRouteCircuitPeeringsClient.Get`

```go
ctx := context.TODO()
id := expressroutecircuitpeerings.NewExpressRouteCircuitPeeringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRouteCircuitValue", "peeringValue")

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
id := expressroutecircuitpeerings.NewExpressRouteCircuitID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRouteCircuitValue")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

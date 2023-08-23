
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/expressroutecircuitconnections` Documentation

The `expressroutecircuitconnections` SDK allows for interaction with the Azure Resource Manager Service `network` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/expressroutecircuitconnections"
```


### Client Initialization

```go
client := expressroutecircuitconnections.NewExpressRouteCircuitConnectionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ExpressRouteCircuitConnectionsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := expressroutecircuitconnections.NewPeeringConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRouteCircuitValue", "peeringValue", "connectionValue")

payload := expressroutecircuitconnections.ExpressRouteCircuitConnection{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ExpressRouteCircuitConnectionsClient.Delete`

```go
ctx := context.TODO()
id := expressroutecircuitconnections.NewPeeringConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRouteCircuitValue", "peeringValue", "connectionValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ExpressRouteCircuitConnectionsClient.Get`

```go
ctx := context.TODO()
id := expressroutecircuitconnections.NewPeeringConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRouteCircuitValue", "peeringValue", "connectionValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ExpressRouteCircuitConnectionsClient.List`

```go
ctx := context.TODO()
id := expressroutecircuitconnections.NewExpressRouteCircuitPeeringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRouteCircuitValue", "peeringValue")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

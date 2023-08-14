
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/expressroutecircuitstats` Documentation

The `expressroutecircuitstats` SDK allows for interaction with the Azure Resource Manager Service `network` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/expressroutecircuitstats"
```


### Client Initialization

```go
client := expressroutecircuitstats.NewExpressRouteCircuitStatsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ExpressRouteCircuitStatsClient.ExpressRouteCircuitsGetPeeringStats`

```go
ctx := context.TODO()
id := expressroutecircuitstats.NewExpressRouteCircuitPeeringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRouteCircuitValue", "peeringValue")

read, err := client.ExpressRouteCircuitsGetPeeringStats(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ExpressRouteCircuitStatsClient.ExpressRouteCircuitsGetStats`

```go
ctx := context.TODO()
id := expressroutecircuitstats.NewExpressRouteCircuitID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRouteCircuitValue")

read, err := client.ExpressRouteCircuitsGetStats(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/expressroutecircuitarptable` Documentation

The `expressroutecircuitarptable` SDK allows for interaction with the Azure Resource Manager Service `network` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/expressroutecircuitarptable"
```


### Client Initialization

```go
client := expressroutecircuitarptable.NewExpressRouteCircuitArpTableClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ExpressRouteCircuitArpTableClient.ExpressRouteCircuitsListArpTable`

```go
ctx := context.TODO()
id := expressroutecircuitarptable.NewArpTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRouteCircuitValue", "peeringValue", "arpTableValue")

if err := client.ExpressRouteCircuitsListArpTableThenPoll(ctx, id); err != nil {
	// handle the error
}
```

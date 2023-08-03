
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/expressroutecircuitroutestablesummary` Documentation

The `expressroutecircuitroutestablesummary` SDK allows for interaction with the Azure Resource Manager Service `network` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/expressroutecircuitroutestablesummary"
```


### Client Initialization

```go
client := expressroutecircuitroutestablesummary.NewExpressRouteCircuitRoutesTableSummaryClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ExpressRouteCircuitRoutesTableSummaryClient.ExpressRouteCircuitsListRoutesTableSummary`

```go
ctx := context.TODO()
id := expressroutecircuitroutestablesummary.NewRouteTablesSummaryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRouteCircuitValue", "peeringValue", "routeTablesSummaryValue")

if err := client.ExpressRouteCircuitsListRoutesTableSummaryThenPoll(ctx, id); err != nil {
	// handle the error
}
```

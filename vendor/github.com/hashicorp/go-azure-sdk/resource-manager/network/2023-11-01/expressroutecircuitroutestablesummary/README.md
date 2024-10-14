
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/expressroutecircuitroutestablesummary` Documentation

The `expressroutecircuitroutestablesummary` SDK allows for interaction with Azure Resource Manager `network` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/expressroutecircuitroutestablesummary"
```


### Client Initialization

```go
client := expressroutecircuitroutestablesummary.NewExpressRouteCircuitRoutesTableSummaryClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ExpressRouteCircuitRoutesTableSummaryClient.ExpressRouteCircuitsListRoutesTableSummary`

```go
ctx := context.TODO()
id := expressroutecircuitroutestablesummary.NewRouteTablesSummaryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRouteCircuitName", "peeringName", "routeTablesSummaryName")

// alternatively `client.ExpressRouteCircuitsListRoutesTableSummary(ctx, id)` can be used to do batched pagination
items, err := client.ExpressRouteCircuitsListRoutesTableSummaryComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

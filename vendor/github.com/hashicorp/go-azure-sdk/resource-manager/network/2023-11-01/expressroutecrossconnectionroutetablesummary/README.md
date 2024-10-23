
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/expressroutecrossconnectionroutetablesummary` Documentation

The `expressroutecrossconnectionroutetablesummary` SDK allows for interaction with Azure Resource Manager `network` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/expressroutecrossconnectionroutetablesummary"
```


### Client Initialization

```go
client := expressroutecrossconnectionroutetablesummary.NewExpressRouteCrossConnectionRouteTableSummaryClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ExpressRouteCrossConnectionRouteTableSummaryClient.ExpressRouteCrossConnectionsListRoutesTableSummary`

```go
ctx := context.TODO()
id := expressroutecrossconnectionroutetablesummary.NewPeeringRouteTablesSummaryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRouteCrossConnectionName", "peeringName", "routeTablesSummaryName")

// alternatively `client.ExpressRouteCrossConnectionsListRoutesTableSummary(ctx, id)` can be used to do batched pagination
items, err := client.ExpressRouteCrossConnectionsListRoutesTableSummaryComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

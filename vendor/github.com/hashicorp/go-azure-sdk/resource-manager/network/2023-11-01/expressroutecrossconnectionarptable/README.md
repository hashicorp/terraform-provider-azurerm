
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/expressroutecrossconnectionarptable` Documentation

The `expressroutecrossconnectionarptable` SDK allows for interaction with Azure Resource Manager `network` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/expressroutecrossconnectionarptable"
```


### Client Initialization

```go
client := expressroutecrossconnectionarptable.NewExpressRouteCrossConnectionArpTableClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ExpressRouteCrossConnectionArpTableClient.ExpressRouteCrossConnectionsListArpTable`

```go
ctx := context.TODO()
id := expressroutecrossconnectionarptable.NewPeeringArpTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRouteCrossConnectionName", "peeringName", "arpTableName")

// alternatively `client.ExpressRouteCrossConnectionsListArpTable(ctx, id)` can be used to do batched pagination
items, err := client.ExpressRouteCrossConnectionsListArpTableComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

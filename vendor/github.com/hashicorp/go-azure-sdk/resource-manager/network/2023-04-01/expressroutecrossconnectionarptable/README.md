
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/expressroutecrossconnectionarptable` Documentation

The `expressroutecrossconnectionarptable` SDK allows for interaction with the Azure Resource Manager Service `network` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/expressroutecrossconnectionarptable"
```


### Client Initialization

```go
client := expressroutecrossconnectionarptable.NewExpressRouteCrossConnectionArpTableClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ExpressRouteCrossConnectionArpTableClient.ExpressRouteCrossConnectionsListArpTable`

```go
ctx := context.TODO()
id := expressroutecrossconnectionarptable.NewPeeringArpTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRouteCrossConnectionValue", "peeringValue", "arpTableValue")

if err := client.ExpressRouteCrossConnectionsListArpTableThenPoll(ctx, id); err != nil {
	// handle the error
}
```

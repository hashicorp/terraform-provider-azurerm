
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/expressroutecrossconnectionpeerings` Documentation

The `expressroutecrossconnectionpeerings` SDK allows for interaction with Azure Resource Manager `network` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/expressroutecrossconnectionpeerings"
```


### Client Initialization

```go
client := expressroutecrossconnectionpeerings.NewExpressRouteCrossConnectionPeeringsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ExpressRouteCrossConnectionPeeringsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := expressroutecrossconnectionpeerings.NewPeeringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRouteCrossConnectionName", "peeringName")

payload := expressroutecrossconnectionpeerings.ExpressRouteCrossConnectionPeering{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ExpressRouteCrossConnectionPeeringsClient.Delete`

```go
ctx := context.TODO()
id := expressroutecrossconnectionpeerings.NewPeeringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRouteCrossConnectionName", "peeringName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ExpressRouteCrossConnectionPeeringsClient.Get`

```go
ctx := context.TODO()
id := expressroutecrossconnectionpeerings.NewPeeringID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRouteCrossConnectionName", "peeringName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ExpressRouteCrossConnectionPeeringsClient.List`

```go
ctx := context.TODO()
id := expressroutecrossconnectionpeerings.NewExpressRouteCrossConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRouteCrossConnectionName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/expressrouteconnections` Documentation

The `expressrouteconnections` SDK allows for interaction with Azure Resource Manager `network` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/expressrouteconnections"
```


### Client Initialization

```go
client := expressrouteconnections.NewExpressRouteConnectionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ExpressRouteConnectionsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := expressrouteconnections.NewExpressRouteConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRouteGatewayName", "expressRouteConnectionName")

payload := expressrouteconnections.ExpressRouteConnection{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ExpressRouteConnectionsClient.Delete`

```go
ctx := context.TODO()
id := expressrouteconnections.NewExpressRouteConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRouteGatewayName", "expressRouteConnectionName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ExpressRouteConnectionsClient.Get`

```go
ctx := context.TODO()
id := expressrouteconnections.NewExpressRouteConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRouteGatewayName", "expressRouteConnectionName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ExpressRouteConnectionsClient.List`

```go
ctx := context.TODO()
id := expressrouteconnections.NewExpressRouteGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRouteGatewayName")

read, err := client.List(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

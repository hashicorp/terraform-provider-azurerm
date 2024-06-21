
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/expressroutecircuitauthorizations` Documentation

The `expressroutecircuitauthorizations` SDK allows for interaction with the Azure Resource Manager Service `network` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/expressroutecircuitauthorizations"
```


### Client Initialization

```go
client := expressroutecircuitauthorizations.NewExpressRouteCircuitAuthorizationsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ExpressRouteCircuitAuthorizationsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := expressroutecircuitauthorizations.NewAuthorizationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRouteCircuitValue", "authorizationValue")

payload := expressroutecircuitauthorizations.ExpressRouteCircuitAuthorization{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ExpressRouteCircuitAuthorizationsClient.Delete`

```go
ctx := context.TODO()
id := expressroutecircuitauthorizations.NewAuthorizationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRouteCircuitValue", "authorizationValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ExpressRouteCircuitAuthorizationsClient.Get`

```go
ctx := context.TODO()
id := expressroutecircuitauthorizations.NewAuthorizationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRouteCircuitValue", "authorizationValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ExpressRouteCircuitAuthorizationsClient.List`

```go
ctx := context.TODO()
id := expressroutecircuitauthorizations.NewExpressRouteCircuitID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRouteCircuitValue")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

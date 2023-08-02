
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/expressroutegateways` Documentation

The `expressroutegateways` SDK allows for interaction with the Azure Resource Manager Service `network` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/expressroutegateways"
```


### Client Initialization

```go
client := expressroutegateways.NewExpressRouteGatewaysClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ExpressRouteGatewaysClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := expressroutegateways.NewExpressRouteGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRouteGatewayValue")

payload := expressroutegateways.ExpressRouteGateway{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ExpressRouteGatewaysClient.Delete`

```go
ctx := context.TODO()
id := expressroutegateways.NewExpressRouteGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRouteGatewayValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ExpressRouteGatewaysClient.Get`

```go
ctx := context.TODO()
id := expressroutegateways.NewExpressRouteGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRouteGatewayValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ExpressRouteGatewaysClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := expressroutegateways.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

read, err := client.ListByResourceGroup(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ExpressRouteGatewaysClient.ListBySubscription`

```go
ctx := context.TODO()
id := expressroutegateways.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

read, err := client.ListBySubscription(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ExpressRouteGatewaysClient.UpdateTags`

```go
ctx := context.TODO()
id := expressroutegateways.NewExpressRouteGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRouteGatewayValue")

payload := expressroutegateways.TagsObject{
	// ...
}


if err := client.UpdateTagsThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/routes` Documentation

The `routes` SDK allows for interaction with the Azure Resource Manager Service `network` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/routes"
```


### Client Initialization

```go
client := routes.NewRoutesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `RoutesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := routes.NewRouteID("12345678-1234-9876-4563-123456789012", "example-resource-group", "routeTableValue", "routeValue")

payload := routes.Route{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `RoutesClient.Delete`

```go
ctx := context.TODO()
id := routes.NewRouteID("12345678-1234-9876-4563-123456789012", "example-resource-group", "routeTableValue", "routeValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `RoutesClient.Get`

```go
ctx := context.TODO()
id := routes.NewRouteID("12345678-1234-9876-4563-123456789012", "example-resource-group", "routeTableValue", "routeValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RoutesClient.List`

```go
ctx := context.TODO()
id := routes.NewRouteTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "routeTableValue")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

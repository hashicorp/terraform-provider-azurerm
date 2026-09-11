
## `github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/routes` Documentation

The `routes` SDK allows for interaction with Azure Resource Manager `cdn` (API Version `2025-12-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/routes"
```


### Client Initialization

```go
client := routes.NewRoutesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `RoutesClient.Create`

```go
ctx := context.TODO()
id := routes.NewRouteID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "afdEndpointName", "routeName")

payload := routes.Route{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `RoutesClient.Delete`

```go
ctx := context.TODO()
id := routes.NewRouteID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "afdEndpointName", "routeName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `RoutesClient.Get`

```go
ctx := context.TODO()
id := routes.NewRouteID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "afdEndpointName", "routeName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RoutesClient.ListByEndpoint`

```go
ctx := context.TODO()
id := routes.NewAfdEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "afdEndpointName")

// alternatively `client.ListByEndpoint(ctx, id)` can be used to do batched pagination
items, err := client.ListByEndpointComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `RoutesClient.Update`

```go
ctx := context.TODO()
id := routes.NewRouteID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "afdEndpointName", "routeName")

payload := routes.RouteUpdateParameters{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

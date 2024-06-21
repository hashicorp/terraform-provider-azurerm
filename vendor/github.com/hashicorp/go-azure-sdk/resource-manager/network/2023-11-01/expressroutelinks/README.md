
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/expressroutelinks` Documentation

The `expressroutelinks` SDK allows for interaction with the Azure Resource Manager Service `network` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/expressroutelinks"
```


### Client Initialization

```go
client := expressroutelinks.NewExpressRouteLinksClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ExpressRouteLinksClient.Get`

```go
ctx := context.TODO()
id := expressroutelinks.NewLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRoutePortValue", "linkValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ExpressRouteLinksClient.List`

```go
ctx := context.TODO()
id := expressroutelinks.NewExpressRoutePortID("12345678-1234-9876-4563-123456789012", "example-resource-group", "expressRoutePortValue")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

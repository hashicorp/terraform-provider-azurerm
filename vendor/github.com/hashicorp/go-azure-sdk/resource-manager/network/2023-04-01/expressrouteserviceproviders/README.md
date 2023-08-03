
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/expressrouteserviceproviders` Documentation

The `expressrouteserviceproviders` SDK allows for interaction with the Azure Resource Manager Service `network` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/expressrouteserviceproviders"
```


### Client Initialization

```go
client := expressrouteserviceproviders.NewExpressRouteServiceProvidersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ExpressRouteServiceProvidersClient.List`

```go
ctx := context.TODO()
id := expressrouteserviceproviders.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

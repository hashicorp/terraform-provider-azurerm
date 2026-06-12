
## `github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/networkconnection` Documentation

The `networkconnection` SDK allows for interaction with Azure Resource Manager `devcenter` (API Version `2025-02-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/networkconnection"
```


### Client Initialization

```go
client := networkconnection.NewNetworkConnectionClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `NetworkConnectionClient.ListOutboundNetworkDependenciesEndpoints`

```go
ctx := context.TODO()
id := networkconnection.NewNetworkConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkConnectionName")

// alternatively `client.ListOutboundNetworkDependenciesEndpoints(ctx, id)` can be used to do batched pagination
items, err := client.ListOutboundNetworkDependenciesEndpointsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NetworkConnectionClient.RunHealthChecks`

```go
ctx := context.TODO()
id := networkconnection.NewNetworkConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkConnectionName")

if err := client.RunHealthChecksThenPoll(ctx, id); err != nil {
	// handle the error
}
```

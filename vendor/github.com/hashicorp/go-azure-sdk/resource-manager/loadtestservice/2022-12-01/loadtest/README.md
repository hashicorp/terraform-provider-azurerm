
## `github.com/hashicorp/go-azure-sdk/resource-manager/loadtestservice/2022-12-01/loadtest` Documentation

The `loadtest` SDK allows for interaction with the Azure Resource Manager Service `loadtestservice` (API Version `2022-12-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/loadtestservice/2022-12-01/loadtest"
```


### Client Initialization

```go
client := loadtest.NewLoadTestClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `LoadTestClient.ListOutboundNetworkDependenciesEndpoints`

```go
ctx := context.TODO()
id := loadtest.NewLoadTestID("12345678-1234-9876-4563-123456789012", "example-resource-group", "loadTestValue")

// alternatively `client.ListOutboundNetworkDependenciesEndpoints(ctx, id)` can be used to do batched pagination
items, err := client.ListOutboundNetworkDependenciesEndpointsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/expressrouteproviderports` Documentation

The `expressrouteproviderports` SDK allows for interaction with the Azure Resource Manager Service `network` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/expressrouteproviderports"
```


### Client Initialization

```go
client := expressrouteproviderports.NewExpressRouteProviderPortsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ExpressRouteProviderPortsClient.ExpressRouteProviderPort`

```go
ctx := context.TODO()
id := expressrouteproviderports.NewExpressRouteProviderPortID("12345678-1234-9876-4563-123456789012", "expressRouteProviderPortValue")

read, err := client.ExpressRouteProviderPort(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ExpressRouteProviderPortsClient.LocationList`

```go
ctx := context.TODO()
id := expressrouteproviderports.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

read, err := client.LocationList(ctx, id, expressrouteproviderports.DefaultLocationListOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

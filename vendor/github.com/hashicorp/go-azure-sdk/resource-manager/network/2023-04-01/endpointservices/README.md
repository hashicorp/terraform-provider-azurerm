
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/endpointservices` Documentation

The `endpointservices` SDK allows for interaction with the Azure Resource Manager Service `network` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/endpointservices"
```


### Client Initialization

```go
client := endpointservices.NewEndpointServicesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `EndpointServicesClient.AvailableEndpointServicesList`

```go
ctx := context.TODO()
id := endpointservices.NewLocationID("12345678-1234-9876-4563-123456789012", "locationValue")

// alternatively `client.AvailableEndpointServicesList(ctx, id)` can be used to do batched pagination
items, err := client.AvailableEndpointServicesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

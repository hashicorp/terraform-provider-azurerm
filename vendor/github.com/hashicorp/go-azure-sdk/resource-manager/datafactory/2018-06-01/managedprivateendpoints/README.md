
## `github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/managedprivateendpoints` Documentation

The `managedprivateendpoints` SDK allows for interaction with the Azure Resource Manager Service `datafactory` (API Version `2018-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/managedprivateendpoints"
```


### Client Initialization

```go
client := managedprivateendpoints.NewManagedPrivateEndpointsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ManagedPrivateEndpointsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := managedprivateendpoints.NewManagedPrivateEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "factoryValue", "managedVirtualNetworkValue", "managedPrivateEndpointValue")

payload := managedprivateendpoints.ManagedPrivateEndpointResource{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, managedprivateendpoints.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedPrivateEndpointsClient.Delete`

```go
ctx := context.TODO()
id := managedprivateendpoints.NewManagedPrivateEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "factoryValue", "managedVirtualNetworkValue", "managedPrivateEndpointValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedPrivateEndpointsClient.Get`

```go
ctx := context.TODO()
id := managedprivateendpoints.NewManagedPrivateEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "factoryValue", "managedVirtualNetworkValue", "managedPrivateEndpointValue")

read, err := client.Get(ctx, id, managedprivateendpoints.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedPrivateEndpointsClient.ListByFactory`

```go
ctx := context.TODO()
id := managedprivateendpoints.NewManagedVirtualNetworkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "factoryValue", "managedVirtualNetworkValue")

// alternatively `client.ListByFactory(ctx, id)` can be used to do batched pagination
items, err := client.ListByFactoryComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

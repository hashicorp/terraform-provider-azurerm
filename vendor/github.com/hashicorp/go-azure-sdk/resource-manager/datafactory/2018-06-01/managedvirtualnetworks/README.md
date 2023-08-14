
## `github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/managedvirtualnetworks` Documentation

The `managedvirtualnetworks` SDK allows for interaction with the Azure Resource Manager Service `datafactory` (API Version `2018-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/managedvirtualnetworks"
```


### Client Initialization

```go
client := managedvirtualnetworks.NewManagedVirtualNetworksClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ManagedVirtualNetworksClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := managedvirtualnetworks.NewManagedVirtualNetworkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "factoryValue", "managedVirtualNetworkValue")

payload := managedvirtualnetworks.ManagedVirtualNetworkResource{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, managedvirtualnetworks.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedVirtualNetworksClient.Get`

```go
ctx := context.TODO()
id := managedvirtualnetworks.NewManagedVirtualNetworkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "factoryValue", "managedVirtualNetworkValue")

read, err := client.Get(ctx, id, managedvirtualnetworks.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedVirtualNetworksClient.ListByFactory`

```go
ctx := context.TODO()
id := managedvirtualnetworks.NewFactoryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "factoryValue")

// alternatively `client.ListByFactory(ctx, id)` can be used to do batched pagination
items, err := client.ListByFactoryComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

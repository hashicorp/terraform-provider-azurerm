
## `github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/virtualnetworks` Documentation

The `virtualnetworks` SDK allows for interaction with the Azure Resource Manager Service `devtestlab` (API Version `2018-09-15`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/virtualnetworks"
```


### Client Initialization

```go
client := virtualnetworks.NewVirtualNetworksClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `VirtualNetworksClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := virtualnetworks.NewVirtualNetworkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "labValue", "virtualNetworkValue")

payload := virtualnetworks.VirtualNetwork{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualNetworksClient.Delete`

```go
ctx := context.TODO()
id := virtualnetworks.NewVirtualNetworkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "labValue", "virtualNetworkValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualNetworksClient.Get`

```go
ctx := context.TODO()
id := virtualnetworks.NewVirtualNetworkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "labValue", "virtualNetworkValue")

read, err := client.Get(ctx, id, virtualnetworks.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualNetworksClient.List`

```go
ctx := context.TODO()
id := virtualnetworks.NewLabID("12345678-1234-9876-4563-123456789012", "example-resource-group", "labValue")

// alternatively `client.List(ctx, id, virtualnetworks.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, virtualnetworks.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualNetworksClient.Update`

```go
ctx := context.TODO()
id := virtualnetworks.NewVirtualNetworkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "labValue", "virtualNetworkValue")

payload := virtualnetworks.UpdateResource{
	// ...
}


read, err := client.Update(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

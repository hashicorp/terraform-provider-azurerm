
## `github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/virtualnetworkaddresses` Documentation

The `virtualnetworkaddresses` SDK allows for interaction with Azure Resource Manager `oracledatabase` (API Version `2024-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/virtualnetworkaddresses"
```


### Client Initialization

```go
client := virtualnetworkaddresses.NewVirtualNetworkAddressesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `VirtualNetworkAddressesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := virtualnetworkaddresses.NewVirtualNetworkAddressID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cloudVmClusterName", "virtualNetworkAddressName")

payload := virtualnetworkaddresses.VirtualNetworkAddress{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualNetworkAddressesClient.Delete`

```go
ctx := context.TODO()
id := virtualnetworkaddresses.NewVirtualNetworkAddressID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cloudVmClusterName", "virtualNetworkAddressName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualNetworkAddressesClient.Get`

```go
ctx := context.TODO()
id := virtualnetworkaddresses.NewVirtualNetworkAddressID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cloudVmClusterName", "virtualNetworkAddressName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualNetworkAddressesClient.ListByCloudVMCluster`

```go
ctx := context.TODO()
id := virtualnetworkaddresses.NewCloudVMClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cloudVmClusterName")

// alternatively `client.ListByCloudVMCluster(ctx, id)` can be used to do batched pagination
items, err := client.ListByCloudVMClusterComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

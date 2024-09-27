
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/subnets` Documentation

The `subnets` SDK allows for interaction with Azure Resource Manager `network` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/subnets"
```


### Client Initialization

```go
client := subnets.NewSubnetsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SubnetsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := commonids.NewSubnetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualNetworkName", "subnetName")

payload := subnets.Subnet{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `SubnetsClient.Delete`

```go
ctx := context.TODO()
id := commonids.NewSubnetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualNetworkName", "subnetName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `SubnetsClient.Get`

```go
ctx := context.TODO()
id := commonids.NewSubnetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualNetworkName", "subnetName")

read, err := client.Get(ctx, id, subnets.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SubnetsClient.List`

```go
ctx := context.TODO()
id := commonids.NewVirtualNetworkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualNetworkName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

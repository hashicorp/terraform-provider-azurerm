
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/virtualrouters` Documentation

The `virtualrouters` SDK allows for interaction with Azure Resource Manager `network` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/virtualrouters"
```


### Client Initialization

```go
client := virtualrouters.NewVirtualRoutersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `VirtualRoutersClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := virtualrouters.NewVirtualRouterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualRouterName")

payload := virtualrouters.VirtualRouter{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualRoutersClient.Delete`

```go
ctx := context.TODO()
id := virtualrouters.NewVirtualRouterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualRouterName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualRoutersClient.Get`

```go
ctx := context.TODO()
id := virtualrouters.NewVirtualRouterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualRouterName")

read, err := client.Get(ctx, id, virtualrouters.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualRoutersClient.List`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualRoutersClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

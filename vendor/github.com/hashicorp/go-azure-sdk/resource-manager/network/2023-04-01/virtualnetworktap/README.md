
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/virtualnetworktap` Documentation

The `virtualnetworktap` SDK allows for interaction with the Azure Resource Manager Service `network` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/virtualnetworktap"
```


### Client Initialization

```go
client := virtualnetworktap.NewVirtualNetworkTapClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `VirtualNetworkTapClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := virtualnetworktap.NewVirtualNetworkTapID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualNetworkTapValue")

payload := virtualnetworktap.VirtualNetworkTap{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualNetworkTapClient.Delete`

```go
ctx := context.TODO()
id := virtualnetworktap.NewVirtualNetworkTapID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualNetworkTapValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualNetworkTapClient.Get`

```go
ctx := context.TODO()
id := virtualnetworktap.NewVirtualNetworkTapID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualNetworkTapValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualNetworkTapClient.UpdateTags`

```go
ctx := context.TODO()
id := virtualnetworktap.NewVirtualNetworkTapID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualNetworkTapValue")

payload := virtualnetworktap.TagsObject{
	// ...
}


read, err := client.UpdateTags(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

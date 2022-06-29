
## `github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2018-09-01/virtualnetworklinks` Documentation

The `virtualnetworklinks` SDK allows for interaction with the Azure Resource Manager Service `privatedns` (API Version `2018-09-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2018-09-01/virtualnetworklinks"
```


### Client Initialization

```go
client := virtualnetworklinks.NewVirtualNetworkLinksClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
if err != nil {
	// handle the error
}
```


### Example Usage: `VirtualNetworkLinksClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := virtualnetworklinks.NewVirtualNetworkLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateZoneValue", "virtualNetworkLinkValue")

payload := virtualnetworklinks.VirtualNetworkLink{
	// ...
}

future, err := client.CreateOrUpdate(ctx, id, payload, virtualnetworklinks.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if err := future.Poller.PollUntilDone(); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualNetworkLinksClient.Delete`

```go
ctx := context.TODO()
id := virtualnetworklinks.NewVirtualNetworkLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateZoneValue", "virtualNetworkLinkValue")
future, err := client.Delete(ctx, id, virtualnetworklinks.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if err := future.Poller.PollUntilDone(); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualNetworkLinksClient.Get`

```go
ctx := context.TODO()
id := virtualnetworklinks.NewVirtualNetworkLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateZoneValue", "virtualNetworkLinkValue")
read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualNetworkLinksClient.List`

```go
ctx := context.TODO()
id := virtualnetworklinks.NewPrivateDnsZoneID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateZoneValue")
// alternatively `client.List(ctx, id, virtualnetworklinks.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, virtualnetworklinks.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualNetworkLinksClient.Update`

```go
ctx := context.TODO()
id := virtualnetworklinks.NewVirtualNetworkLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateZoneValue", "virtualNetworkLinkValue")

payload := virtualnetworklinks.VirtualNetworkLink{
	// ...
}

future, err := client.Update(ctx, id, payload, virtualnetworklinks.DefaultUpdateOperationOptions())
if err != nil {
	// handle the error
}
if err := future.Poller.PollUntilDone(); err != nil {
	// handle the error
}
```

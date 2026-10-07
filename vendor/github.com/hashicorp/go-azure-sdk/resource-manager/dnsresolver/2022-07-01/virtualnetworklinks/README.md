
## `github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/virtualnetworklinks` Documentation

The `virtualnetworklinks` SDK allows for interaction with Azure Resource Manager `dnsresolver` (API Version `2022-07-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/virtualnetworklinks"
```


### Client Initialization

```go
client := virtualnetworklinks.NewVirtualNetworkLinksClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `VirtualNetworkLinksClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := virtualnetworklinks.NewVirtualNetworkLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dnsForwardingRulesetName", "virtualNetworkLinkName")

payload := virtualnetworklinks.VirtualNetworkLink{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload, virtualnetworklinks.DefaultCreateOrUpdateOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualNetworkLinksClient.Delete`

```go
ctx := context.TODO()
id := virtualnetworklinks.NewVirtualNetworkLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dnsForwardingRulesetName", "virtualNetworkLinkName")

if err := client.DeleteThenPoll(ctx, id, virtualnetworklinks.DefaultDeleteOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualNetworkLinksClient.Get`

```go
ctx := context.TODO()
id := virtualnetworklinks.NewVirtualNetworkLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dnsForwardingRulesetName", "virtualNetworkLinkName")

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
id := virtualnetworklinks.NewDnsForwardingRulesetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dnsForwardingRulesetName")

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
id := virtualnetworklinks.NewVirtualNetworkLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dnsForwardingRulesetName", "virtualNetworkLinkName")

payload := virtualnetworklinks.VirtualNetworkLinkPatch{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload, virtualnetworklinks.DefaultUpdateOperationOptions()); err != nil {
	// handle the error
}
```

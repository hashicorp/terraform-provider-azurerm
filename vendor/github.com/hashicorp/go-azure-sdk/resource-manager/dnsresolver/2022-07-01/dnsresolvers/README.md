
## `github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/dnsresolvers` Documentation

The `dnsresolvers` SDK allows for interaction with the Azure Resource Manager Service `dnsresolver` (API Version `2022-07-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/dnsresolvers"
```


### Client Initialization

```go
client := dnsresolvers.NewDnsResolversClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DnsResolversClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := dnsresolvers.NewDnsResolverID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dnsResolverValue")

payload := dnsresolvers.DnsResolver{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload, dnsresolvers.DefaultCreateOrUpdateOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `DnsResolversClient.Delete`

```go
ctx := context.TODO()
id := dnsresolvers.NewDnsResolverID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dnsResolverValue")

if err := client.DeleteThenPoll(ctx, id, dnsresolvers.DefaultDeleteOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `DnsResolversClient.Get`

```go
ctx := context.TODO()
id := dnsresolvers.NewDnsResolverID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dnsResolverValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DnsResolversClient.List`

```go
ctx := context.TODO()
id := dnsresolvers.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.List(ctx, id, dnsresolvers.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, dnsresolvers.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DnsResolversClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := dnsresolvers.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id, dnsresolvers.DefaultListByResourceGroupOperationOptions())` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id, dnsresolvers.DefaultListByResourceGroupOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DnsResolversClient.ListByVirtualNetwork`

```go
ctx := context.TODO()
id := dnsresolvers.NewVirtualNetworkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualNetworkValue")

// alternatively `client.ListByVirtualNetwork(ctx, id, dnsresolvers.DefaultListByVirtualNetworkOperationOptions())` can be used to do batched pagination
items, err := client.ListByVirtualNetworkComplete(ctx, id, dnsresolvers.DefaultListByVirtualNetworkOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DnsResolversClient.Update`

```go
ctx := context.TODO()
id := dnsresolvers.NewDnsResolverID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dnsResolverValue")

payload := dnsresolvers.DnsResolverPatch{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload, dnsresolvers.DefaultUpdateOperationOptions()); err != nil {
	// handle the error
}
```

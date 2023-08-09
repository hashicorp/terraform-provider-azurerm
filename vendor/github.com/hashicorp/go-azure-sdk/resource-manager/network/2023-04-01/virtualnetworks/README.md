
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/virtualnetworks` Documentation

The `virtualnetworks` SDK allows for interaction with the Azure Resource Manager Service `network` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/virtualnetworks"
```


### Client Initialization

```go
client := virtualnetworks.NewVirtualNetworksClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `VirtualNetworksClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := virtualnetworks.NewVirtualNetworkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualNetworkValue")

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
id := virtualnetworks.NewVirtualNetworkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualNetworkValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualNetworksClient.Get`

```go
ctx := context.TODO()
id := virtualnetworks.NewVirtualNetworkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualNetworkValue")

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
id := virtualnetworks.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualNetworksClient.ListAll`

```go
ctx := context.TODO()
id := virtualnetworks.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListAll(ctx, id)` can be used to do batched pagination
items, err := client.ListAllComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualNetworksClient.ResourceNavigationLinksList`

```go
ctx := context.TODO()
id := virtualnetworks.NewSubnetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualNetworkValue", "subnetValue")

read, err := client.ResourceNavigationLinksList(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualNetworksClient.ServiceAssociationLinksList`

```go
ctx := context.TODO()
id := virtualnetworks.NewSubnetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualNetworkValue", "subnetValue")

read, err := client.ServiceAssociationLinksList(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualNetworksClient.SubnetsPrepareNetworkPolicies`

```go
ctx := context.TODO()
id := virtualnetworks.NewSubnetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualNetworkValue", "subnetValue")

payload := virtualnetworks.PrepareNetworkPoliciesRequest{
	// ...
}


if err := client.SubnetsPrepareNetworkPoliciesThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualNetworksClient.SubnetsUnprepareNetworkPolicies`

```go
ctx := context.TODO()
id := virtualnetworks.NewSubnetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualNetworkValue", "subnetValue")

payload := virtualnetworks.UnprepareNetworkPoliciesRequest{
	// ...
}


if err := client.SubnetsUnprepareNetworkPoliciesThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualNetworksClient.UpdateTags`

```go
ctx := context.TODO()
id := virtualnetworks.NewVirtualNetworkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualNetworkValue")

payload := virtualnetworks.TagsObject{
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


### Example Usage: `VirtualNetworksClient.VirtualNetworksCheckIPAddressAvailability`

```go
ctx := context.TODO()
id := virtualnetworks.NewVirtualNetworkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualNetworkValue")

read, err := client.VirtualNetworksCheckIPAddressAvailability(ctx, id, virtualnetworks.DefaultVirtualNetworksCheckIPAddressAvailabilityOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualNetworksClient.VirtualNetworksListDdosProtectionStatus`

```go
ctx := context.TODO()
id := virtualnetworks.NewVirtualNetworkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualNetworkValue")

// alternatively `client.VirtualNetworksListDdosProtectionStatus(ctx, id, virtualnetworks.DefaultVirtualNetworksListDdosProtectionStatusOperationOptions())` can be used to do batched pagination
items, err := client.VirtualNetworksListDdosProtectionStatusComplete(ctx, id, virtualnetworks.DefaultVirtualNetworksListDdosProtectionStatusOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualNetworksClient.VirtualNetworksListUsage`

```go
ctx := context.TODO()
id := virtualnetworks.NewVirtualNetworkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualNetworkValue")

// alternatively `client.VirtualNetworksListUsage(ctx, id)` can be used to do batched pagination
items, err := client.VirtualNetworksListUsageComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

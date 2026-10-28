
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualnetworks` Documentation

The `virtualnetworks` SDK allows for interaction with Azure Resource Manager `network` (API Version `2025-07-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualnetworks"
```


### Client Initialization

```go
client := virtualnetworks.NewVirtualNetworksClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `VirtualNetworksClient.AvailableDelegationsList`

```go
ctx := context.TODO()
id := virtualnetworks.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

// alternatively `client.AvailableDelegationsList(ctx, id)` can be used to do batched pagination
items, err := client.AvailableDelegationsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualNetworksClient.AvailableEndpointServicesList`

```go
ctx := context.TODO()
id := virtualnetworks.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

// alternatively `client.AvailableEndpointServicesList(ctx, id)` can be used to do batched pagination
items, err := client.AvailableEndpointServicesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualNetworksClient.AvailablePrivateEndpointTypesList`

```go
ctx := context.TODO()
id := virtualnetworks.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

// alternatively `client.AvailablePrivateEndpointTypesList(ctx, id)` can be used to do batched pagination
items, err := client.AvailablePrivateEndpointTypesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualNetworksClient.AvailablePrivateEndpointTypesListByResourceGroup`

```go
ctx := context.TODO()
id := virtualnetworks.NewProviderLocationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "locationName")

// alternatively `client.AvailablePrivateEndpointTypesListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.AvailablePrivateEndpointTypesListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualNetworksClient.AvailableResourceGroupDelegationsList`

```go
ctx := context.TODO()
id := virtualnetworks.NewProviderLocationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "locationName")

// alternatively `client.AvailableResourceGroupDelegationsList(ctx, id)` can be used to do batched pagination
items, err := client.AvailableResourceGroupDelegationsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualNetworksClient.AvailableServiceAliasesList`

```go
ctx := context.TODO()
id := virtualnetworks.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

// alternatively `client.AvailableServiceAliasesList(ctx, id)` can be used to do batched pagination
items, err := client.AvailableServiceAliasesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualNetworksClient.AvailableServiceAliasesListByResourceGroup`

```go
ctx := context.TODO()
id := virtualnetworks.NewProviderLocationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "locationName")

// alternatively `client.AvailableServiceAliasesListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.AvailableServiceAliasesListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualNetworksClient.CheckDnsNameAvailability`

```go
ctx := context.TODO()
id := virtualnetworks.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

read, err := client.CheckDnsNameAvailability(ctx, id, virtualnetworks.DefaultCheckDnsNameAvailabilityOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualNetworksClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := commonids.NewVirtualNetworkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualNetworkName")

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
id := commonids.NewVirtualNetworkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualNetworkName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualNetworksClient.Get`

```go
ctx := context.TODO()
id := commonids.NewVirtualNetworkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualNetworkName")

read, err := client.Get(ctx, id, virtualnetworks.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualNetworksClient.InboundSecurityRuleCreateOrUpdate`

```go
ctx := context.TODO()
id := virtualnetworks.NewInboundSecurityRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkVirtualApplianceName", "inboundSecurityRuleName")

payload := virtualnetworks.InboundSecurityRule{
	// ...
}


if err := client.InboundSecurityRuleCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualNetworksClient.InboundSecurityRuleGet`

```go
ctx := context.TODO()
id := virtualnetworks.NewInboundSecurityRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkVirtualApplianceName", "inboundSecurityRuleName")

read, err := client.InboundSecurityRuleGet(ctx, id)
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
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

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
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListAll(ctx, id)` can be used to do batched pagination
items, err := client.ListAllComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualNetworksClient.PublicIPAddressesGetCloudServicePublicIPAddress`

```go
ctx := context.TODO()
id := commonids.NewCloudServicesPublicIPAddressID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cloudServiceName", "roleInstanceName", "networkInterfaceName", "ipConfigurationName", "publicIPAddressName")

read, err := client.PublicIPAddressesGetCloudServicePublicIPAddress(ctx, id, virtualnetworks.DefaultPublicIPAddressesGetCloudServicePublicIPAddressOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualNetworksClient.PublicIPAddressesListCloudServicePublicIPAddresses`

```go
ctx := context.TODO()
id := virtualnetworks.NewProviderCloudServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cloudServiceName")

// alternatively `client.PublicIPAddressesListCloudServicePublicIPAddresses(ctx, id)` can be used to do batched pagination
items, err := client.PublicIPAddressesListCloudServicePublicIPAddressesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualNetworksClient.PublicIPAddressesListCloudServiceRoleInstancePublicIPAddresses`

```go
ctx := context.TODO()
id := commonids.NewCloudServicesIPConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cloudServiceName", "roleInstanceName", "networkInterfaceName", "ipConfigurationName")

// alternatively `client.PublicIPAddressesListCloudServiceRoleInstancePublicIPAddresses(ctx, id)` can be used to do batched pagination
items, err := client.PublicIPAddressesListCloudServiceRoleInstancePublicIPAddressesComplete(ctx, id)
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
id := commonids.NewSubnetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualNetworkName", "subnetName")

// alternatively `client.ResourceNavigationLinksList(ctx, id)` can be used to do batched pagination
items, err := client.ResourceNavigationLinksListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualNetworksClient.ServiceAssociationLinksList`

```go
ctx := context.TODO()
id := commonids.NewSubnetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualNetworkName", "subnetName")

// alternatively `client.ServiceAssociationLinksList(ctx, id)` can be used to do batched pagination
items, err := client.ServiceAssociationLinksListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualNetworksClient.ServiceTagInformationList`

```go
ctx := context.TODO()
id := virtualnetworks.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

// alternatively `client.ServiceTagInformationList(ctx, id, virtualnetworks.DefaultServiceTagInformationListOperationOptions())` can be used to do batched pagination
items, err := client.ServiceTagInformationListComplete(ctx, id, virtualnetworks.DefaultServiceTagInformationListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualNetworksClient.ServiceTagsList`

```go
ctx := context.TODO()
id := virtualnetworks.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

read, err := client.ServiceTagsList(ctx, id)
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
id := commonids.NewSubnetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualNetworkName", "subnetName")

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
id := commonids.NewSubnetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualNetworkName", "subnetName")

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
id := commonids.NewVirtualNetworkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualNetworkName")

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
id := commonids.NewVirtualNetworkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualNetworkName")

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
id := commonids.NewVirtualNetworkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualNetworkName")

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
id := commonids.NewVirtualNetworkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualNetworkName")

// alternatively `client.VirtualNetworksListUsage(ctx, id)` can be used to do batched pagination
items, err := client.VirtualNetworksListUsageComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

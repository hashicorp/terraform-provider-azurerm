
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/loadbalancers` Documentation

The `loadbalancers` SDK allows for interaction with Azure Resource Manager `network` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/loadbalancers"
```


### Client Initialization

```go
client := loadbalancers.NewLoadBalancersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `LoadBalancersClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := loadbalancers.NewProviderLoadBalancerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "loadBalancerName")

payload := loadbalancers.LoadBalancer{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `LoadBalancersClient.Delete`

```go
ctx := context.TODO()
id := loadbalancers.NewProviderLoadBalancerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "loadBalancerName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `LoadBalancersClient.Get`

```go
ctx := context.TODO()
id := loadbalancers.NewProviderLoadBalancerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "loadBalancerName")

read, err := client.Get(ctx, id, loadbalancers.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LoadBalancersClient.InboundNatRulesCreateOrUpdate`

```go
ctx := context.TODO()
id := loadbalancers.NewInboundNatRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "loadBalancerName", "inboundNatRuleName")

payload := loadbalancers.InboundNatRule{
	// ...
}


if err := client.InboundNatRulesCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `LoadBalancersClient.InboundNatRulesDelete`

```go
ctx := context.TODO()
id := loadbalancers.NewInboundNatRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "loadBalancerName", "inboundNatRuleName")

if err := client.InboundNatRulesDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `LoadBalancersClient.InboundNatRulesGet`

```go
ctx := context.TODO()
id := loadbalancers.NewInboundNatRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "loadBalancerName", "inboundNatRuleName")

read, err := client.InboundNatRulesGet(ctx, id, loadbalancers.DefaultInboundNatRulesGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LoadBalancersClient.InboundNatRulesList`

```go
ctx := context.TODO()
id := loadbalancers.NewProviderLoadBalancerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "loadBalancerName")

// alternatively `client.InboundNatRulesList(ctx, id)` can be used to do batched pagination
items, err := client.InboundNatRulesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `LoadBalancersClient.List`

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


### Example Usage: `LoadBalancersClient.ListAll`

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


### Example Usage: `LoadBalancersClient.ListInboundNatRulePortMappings`

```go
ctx := context.TODO()
id := loadbalancers.NewBackendAddressPoolID("12345678-1234-9876-4563-123456789012", "resourceGroupName", "loadBalancerName", "backendAddressPoolName")

payload := loadbalancers.QueryInboundNatRulePortMappingRequest{
	// ...
}


if err := client.ListInboundNatRulePortMappingsThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `LoadBalancersClient.LoadBalancerBackendAddressPoolsCreateOrUpdate`

```go
ctx := context.TODO()
id := loadbalancers.NewLoadBalancerBackendAddressPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "loadBalancerName", "backendAddressPoolName")

payload := loadbalancers.BackendAddressPool{
	// ...
}


if err := client.LoadBalancerBackendAddressPoolsCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `LoadBalancersClient.LoadBalancerBackendAddressPoolsDelete`

```go
ctx := context.TODO()
id := loadbalancers.NewLoadBalancerBackendAddressPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "loadBalancerName", "backendAddressPoolName")

if err := client.LoadBalancerBackendAddressPoolsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `LoadBalancersClient.LoadBalancerBackendAddressPoolsGet`

```go
ctx := context.TODO()
id := loadbalancers.NewLoadBalancerBackendAddressPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "loadBalancerName", "backendAddressPoolName")

read, err := client.LoadBalancerBackendAddressPoolsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LoadBalancersClient.LoadBalancerBackendAddressPoolsList`

```go
ctx := context.TODO()
id := loadbalancers.NewProviderLoadBalancerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "loadBalancerName")

// alternatively `client.LoadBalancerBackendAddressPoolsList(ctx, id)` can be used to do batched pagination
items, err := client.LoadBalancerBackendAddressPoolsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `LoadBalancersClient.LoadBalancerFrontendIPConfigurationsGet`

```go
ctx := context.TODO()
id := loadbalancers.NewFrontendIPConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "loadBalancerName", "frontendIPConfigurationName")

read, err := client.LoadBalancerFrontendIPConfigurationsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LoadBalancersClient.LoadBalancerFrontendIPConfigurationsList`

```go
ctx := context.TODO()
id := loadbalancers.NewProviderLoadBalancerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "loadBalancerName")

// alternatively `client.LoadBalancerFrontendIPConfigurationsList(ctx, id)` can be used to do batched pagination
items, err := client.LoadBalancerFrontendIPConfigurationsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `LoadBalancersClient.LoadBalancerLoadBalancingRulesGet`

```go
ctx := context.TODO()
id := loadbalancers.NewLoadBalancingRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "loadBalancerName", "loadBalancingRuleName")

read, err := client.LoadBalancerLoadBalancingRulesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LoadBalancersClient.LoadBalancerLoadBalancingRulesList`

```go
ctx := context.TODO()
id := loadbalancers.NewProviderLoadBalancerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "loadBalancerName")

// alternatively `client.LoadBalancerLoadBalancingRulesList(ctx, id)` can be used to do batched pagination
items, err := client.LoadBalancerLoadBalancingRulesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `LoadBalancersClient.LoadBalancerNetworkInterfacesList`

```go
ctx := context.TODO()
id := loadbalancers.NewProviderLoadBalancerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "loadBalancerName")

// alternatively `client.LoadBalancerNetworkInterfacesList(ctx, id)` can be used to do batched pagination
items, err := client.LoadBalancerNetworkInterfacesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `LoadBalancersClient.LoadBalancerOutboundRulesGet`

```go
ctx := context.TODO()
id := loadbalancers.NewOutboundRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "loadBalancerName", "outboundRuleName")

read, err := client.LoadBalancerOutboundRulesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LoadBalancersClient.LoadBalancerOutboundRulesList`

```go
ctx := context.TODO()
id := loadbalancers.NewProviderLoadBalancerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "loadBalancerName")

// alternatively `client.LoadBalancerOutboundRulesList(ctx, id)` can be used to do batched pagination
items, err := client.LoadBalancerOutboundRulesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `LoadBalancersClient.LoadBalancerProbesGet`

```go
ctx := context.TODO()
id := loadbalancers.NewProbeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "loadBalancerName", "probeName")

read, err := client.LoadBalancerProbesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LoadBalancersClient.LoadBalancerProbesList`

```go
ctx := context.TODO()
id := loadbalancers.NewProviderLoadBalancerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "loadBalancerName")

// alternatively `client.LoadBalancerProbesList(ctx, id)` can be used to do batched pagination
items, err := client.LoadBalancerProbesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `LoadBalancersClient.MigrateToIPBased`

```go
ctx := context.TODO()
id := loadbalancers.NewLoadBalancerID("12345678-1234-9876-4563-123456789012", "resourceGroupName", "loadBalancerName")

payload := loadbalancers.MigrateLoadBalancerToIPBasedRequest{
	// ...
}


read, err := client.MigrateToIPBased(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LoadBalancersClient.SwapPublicIPAddresses`

```go
ctx := context.TODO()
id := loadbalancers.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

payload := loadbalancers.LoadBalancerVipSwapRequest{
	// ...
}


if err := client.SwapPublicIPAddressesThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `LoadBalancersClient.UpdateTags`

```go
ctx := context.TODO()
id := loadbalancers.NewProviderLoadBalancerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "loadBalancerName")

payload := loadbalancers.TagsObject{
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

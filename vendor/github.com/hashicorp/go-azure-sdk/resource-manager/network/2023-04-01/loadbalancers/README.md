
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/loadbalancers` Documentation

The `loadbalancers` SDK allows for interaction with the Azure Resource Manager Service `network` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/loadbalancers"
```


### Client Initialization

```go
client := loadbalancers.NewLoadBalancersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `LoadBalancersClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := loadbalancers.NewProviderLoadBalancerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "loadBalancerValue")

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
id := loadbalancers.NewProviderLoadBalancerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "loadBalancerValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `LoadBalancersClient.Get`

```go
ctx := context.TODO()
id := loadbalancers.NewProviderLoadBalancerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "loadBalancerValue")

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
id := loadbalancers.NewInboundNatRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "loadBalancerValue", "inboundNatRuleValue")

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
id := loadbalancers.NewInboundNatRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "loadBalancerValue", "inboundNatRuleValue")

if err := client.InboundNatRulesDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `LoadBalancersClient.InboundNatRulesGet`

```go
ctx := context.TODO()
id := loadbalancers.NewInboundNatRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "loadBalancerValue", "inboundNatRuleValue")

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
id := loadbalancers.NewProviderLoadBalancerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "loadBalancerValue")

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
id := loadbalancers.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

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
id := loadbalancers.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

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
id := loadbalancers.NewBackendAddressPoolID("12345678-1234-9876-4563-123456789012", "resourceGroupValue", "loadBalancerValue", "backendAddressPoolValue")

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
id := loadbalancers.NewLoadBalancerBackendAddressPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "loadBalancerValue", "backendAddressPoolValue")

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
id := loadbalancers.NewLoadBalancerBackendAddressPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "loadBalancerValue", "backendAddressPoolValue")

if err := client.LoadBalancerBackendAddressPoolsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `LoadBalancersClient.LoadBalancerBackendAddressPoolsGet`

```go
ctx := context.TODO()
id := loadbalancers.NewLoadBalancerBackendAddressPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "loadBalancerValue", "backendAddressPoolValue")

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
id := loadbalancers.NewProviderLoadBalancerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "loadBalancerValue")

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
id := loadbalancers.NewFrontendIPConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "loadBalancerValue", "frontendIPConfigurationValue")

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
id := loadbalancers.NewProviderLoadBalancerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "loadBalancerValue")

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
id := loadbalancers.NewLoadBalancingRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "loadBalancerValue", "loadBalancingRuleValue")

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
id := loadbalancers.NewProviderLoadBalancerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "loadBalancerValue")

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
id := loadbalancers.NewProviderLoadBalancerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "loadBalancerValue")

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
id := loadbalancers.NewOutboundRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "loadBalancerValue", "outboundRuleValue")

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
id := loadbalancers.NewProviderLoadBalancerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "loadBalancerValue")

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
id := loadbalancers.NewProbeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "loadBalancerValue", "probeValue")

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
id := loadbalancers.NewProviderLoadBalancerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "loadBalancerValue")

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
id := loadbalancers.NewLoadBalancerID("12345678-1234-9876-4563-123456789012", "resourceGroupValue", "loadBalancerValue")

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
id := loadbalancers.NewLocationID("12345678-1234-9876-4563-123456789012", "locationValue")

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
id := loadbalancers.NewProviderLoadBalancerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "loadBalancerValue")

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

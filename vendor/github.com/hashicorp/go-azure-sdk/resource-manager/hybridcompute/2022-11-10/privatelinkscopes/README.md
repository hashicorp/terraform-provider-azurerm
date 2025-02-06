
## `github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2022-11-10/privatelinkscopes` Documentation

The `privatelinkscopes` SDK allows for interaction with Azure Resource Manager `hybridcompute` (API Version `2022-11-10`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2022-11-10/privatelinkscopes"
```


### Client Initialization

```go
client := privatelinkscopes.NewPrivateLinkScopesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PrivateLinkScopesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := privatelinkscopes.NewProviderPrivateLinkScopeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateLinkScopeName")

payload := privatelinkscopes.HybridComputePrivateLinkScope{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PrivateLinkScopesClient.Delete`

```go
ctx := context.TODO()
id := privatelinkscopes.NewProviderPrivateLinkScopeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateLinkScopeName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `PrivateLinkScopesClient.Get`

```go
ctx := context.TODO()
id := privatelinkscopes.NewProviderPrivateLinkScopeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateLinkScopeName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PrivateLinkScopesClient.GetValidationDetails`

```go
ctx := context.TODO()
id := privatelinkscopes.NewPrivateLinkScopeID("12345678-1234-9876-4563-123456789012", "locationName", "privateLinkScopeId")

read, err := client.GetValidationDetails(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PrivateLinkScopesClient.GetValidationDetailsForMachine`

```go
ctx := context.TODO()
id := privatelinkscopes.NewMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "machineName")

read, err := client.GetValidationDetailsForMachine(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PrivateLinkScopesClient.List`

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


### Example Usage: `PrivateLinkScopesClient.ListByResourceGroup`

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


### Example Usage: `PrivateLinkScopesClient.UpdateTags`

```go
ctx := context.TODO()
id := privatelinkscopes.NewProviderPrivateLinkScopeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateLinkScopeName")

payload := privatelinkscopes.TagsResource{
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

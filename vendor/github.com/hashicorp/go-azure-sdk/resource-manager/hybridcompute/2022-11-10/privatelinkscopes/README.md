
## `github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2022-11-10/privatelinkscopes` Documentation

The `privatelinkscopes` SDK allows for interaction with the Azure Resource Manager Service `hybridcompute` (API Version `2022-11-10`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2022-11-10/privatelinkscopes"
```


### Client Initialization

```go
client := privatelinkscopes.NewPrivateLinkScopesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PrivateLinkScopesClient.PrivateLinkScopesCreateOrUpdate`

```go
ctx := context.TODO()
id := privatelinkscopes.NewProviderPrivateLinkScopeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateLinkScopeValue")

payload := privatelinkscopes.HybridComputePrivateLinkScope{
	// ...
}


read, err := client.PrivateLinkScopesCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PrivateLinkScopesClient.PrivateLinkScopesDelete`

```go
ctx := context.TODO()
id := privatelinkscopes.NewProviderPrivateLinkScopeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateLinkScopeValue")

if err := client.PrivateLinkScopesDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `PrivateLinkScopesClient.PrivateLinkScopesGet`

```go
ctx := context.TODO()
id := privatelinkscopes.NewProviderPrivateLinkScopeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateLinkScopeValue")

read, err := client.PrivateLinkScopesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PrivateLinkScopesClient.PrivateLinkScopesGetValidationDetails`

```go
ctx := context.TODO()
id := privatelinkscopes.NewPrivateLinkScopeID("12345678-1234-9876-4563-123456789012", "locationValue", "privateLinkScopeIdValue")

read, err := client.PrivateLinkScopesGetValidationDetails(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PrivateLinkScopesClient.PrivateLinkScopesGetValidationDetailsForMachine`

```go
ctx := context.TODO()
id := privatelinkscopes.NewMachineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "machineValue")

read, err := client.PrivateLinkScopesGetValidationDetailsForMachine(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PrivateLinkScopesClient.PrivateLinkScopesList`

```go
ctx := context.TODO()
id := privatelinkscopes.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.PrivateLinkScopesList(ctx, id)` can be used to do batched pagination
items, err := client.PrivateLinkScopesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PrivateLinkScopesClient.PrivateLinkScopesListByResourceGroup`

```go
ctx := context.TODO()
id := privatelinkscopes.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.PrivateLinkScopesListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.PrivateLinkScopesListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PrivateLinkScopesClient.PrivateLinkScopesUpdateTags`

```go
ctx := context.TODO()
id := privatelinkscopes.NewProviderPrivateLinkScopeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateLinkScopeValue")

payload := privatelinkscopes.TagsResource{
	// ...
}


read, err := client.PrivateLinkScopesUpdateTags(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

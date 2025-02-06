
## `github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-09-01/providers` Documentation

The `providers` SDK allows for interaction with Azure Resource Manager `resources` (API Version `2022-09-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-09-01/providers"
```


### Client Initialization

```go
client := providers.NewProvidersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ProvidersClient.Get`

```go
ctx := context.TODO()
id := providers.NewSubscriptionProviderID("12345678-1234-9876-4563-123456789012", "providerName")

read, err := client.Get(ctx, id, providers.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProvidersClient.GetAtTenantScope`

```go
ctx := context.TODO()
id := providers.NewProviderID("providerName")

read, err := client.GetAtTenantScope(ctx, id, providers.DefaultGetAtTenantScopeOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProvidersClient.List`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.List(ctx, id, providers.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, providers.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ProvidersClient.ListAtTenantScope`

```go
ctx := context.TODO()


// alternatively `client.ListAtTenantScope(ctx, providers.DefaultListAtTenantScopeOperationOptions())` can be used to do batched pagination
items, err := client.ListAtTenantScopeComplete(ctx, providers.DefaultListAtTenantScopeOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ProvidersClient.ProviderPermissions`

```go
ctx := context.TODO()
id := providers.NewSubscriptionProviderID("12345678-1234-9876-4563-123456789012", "providerName")

// alternatively `client.ProviderPermissions(ctx, id)` can be used to do batched pagination
items, err := client.ProviderPermissionsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ProvidersClient.ProviderResourceTypesList`

```go
ctx := context.TODO()
id := providers.NewSubscriptionProviderID("12345678-1234-9876-4563-123456789012", "providerName")

// alternatively `client.ProviderResourceTypesList(ctx, id, providers.DefaultProviderResourceTypesListOperationOptions())` can be used to do batched pagination
items, err := client.ProviderResourceTypesListComplete(ctx, id, providers.DefaultProviderResourceTypesListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ProvidersClient.Register`

```go
ctx := context.TODO()
id := providers.NewSubscriptionProviderID("12345678-1234-9876-4563-123456789012", "providerName")

payload := providers.ProviderRegistrationRequest{
	// ...
}


read, err := client.Register(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProvidersClient.RegisterAtManagementGroupScope`

```go
ctx := context.TODO()
id := providers.NewProviders2ID("groupId", "providerName")

read, err := client.RegisterAtManagementGroupScope(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProvidersClient.Unregister`

```go
ctx := context.TODO()
id := providers.NewSubscriptionProviderID("12345678-1234-9876-4563-123456789012", "providerName")

read, err := client.Unregister(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

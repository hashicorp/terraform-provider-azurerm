
## `github.com/hashicorp/go-azure-sdk/resource-manager/managedidentity/2018-11-30/managedidentities` Documentation

The `managedidentities` SDK allows for interaction with the Azure Resource Manager Service `managedidentity` (API Version `2018-11-30`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/managedidentity/2018-11-30/managedidentities"
```


### Client Initialization

```go
client := managedidentities.NewManagedIdentitiesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ManagedIdentitiesClient.SystemAssignedIdentitiesGetByScope`

```go
ctx := context.TODO()
id := managedidentities.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

read, err := client.SystemAssignedIdentitiesGetByScope(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedIdentitiesClient.UserAssignedIdentitiesCreateOrUpdate`

```go
ctx := context.TODO()
id := managedidentities.NewUserAssignedIdentityID("12345678-1234-9876-4563-123456789012", "example-resource-group", "resourceValue")

payload := managedidentities.Identity{
	// ...
}


read, err := client.UserAssignedIdentitiesCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedIdentitiesClient.UserAssignedIdentitiesDelete`

```go
ctx := context.TODO()
id := managedidentities.NewUserAssignedIdentityID("12345678-1234-9876-4563-123456789012", "example-resource-group", "resourceValue")

read, err := client.UserAssignedIdentitiesDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedIdentitiesClient.UserAssignedIdentitiesGet`

```go
ctx := context.TODO()
id := managedidentities.NewUserAssignedIdentityID("12345678-1234-9876-4563-123456789012", "example-resource-group", "resourceValue")

read, err := client.UserAssignedIdentitiesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedIdentitiesClient.UserAssignedIdentitiesListByResourceGroup`

```go
ctx := context.TODO()
id := managedidentities.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.UserAssignedIdentitiesListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.UserAssignedIdentitiesListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ManagedIdentitiesClient.UserAssignedIdentitiesListBySubscription`

```go
ctx := context.TODO()
id := managedidentities.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.UserAssignedIdentitiesListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.UserAssignedIdentitiesListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ManagedIdentitiesClient.UserAssignedIdentitiesUpdate`

```go
ctx := context.TODO()
id := managedidentities.NewUserAssignedIdentityID("12345678-1234-9876-4563-123456789012", "example-resource-group", "resourceValue")

payload := managedidentities.IdentityUpdate{
	// ...
}


read, err := client.UserAssignedIdentitiesUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

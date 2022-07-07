
## `github.com/hashicorp/go-azure-sdk/resource-manager/managedidentity/2018-11-30/managedidentity` Documentation

The `managedidentity` SDK allows for interaction with the Azure Resource Manager Service `managedidentity` (API Version `2018-11-30`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/managedidentity/2018-11-30/managedidentity"
```


### Client Initialization

```go
client := managedidentity.NewManagedIdentityClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ManagedIdentityClient.SystemAssignedIdentitiesGetByScope`

```go
ctx := context.TODO()
id := managedidentity.NewScopeID()

read, err := client.SystemAssignedIdentitiesGetByScope(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedIdentityClient.UserAssignedIdentitiesCreateOrUpdate`

```go
ctx := context.TODO()
id := managedidentity.NewUserAssignedIdentityID()

payload := managedidentity.Identity{
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


### Example Usage: `ManagedIdentityClient.UserAssignedIdentitiesDelete`

```go
ctx := context.TODO()
id := managedidentity.NewUserAssignedIdentityID()

read, err := client.UserAssignedIdentitiesDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedIdentityClient.UserAssignedIdentitiesGet`

```go
ctx := context.TODO()
id := managedidentity.NewUserAssignedIdentityID()

read, err := client.UserAssignedIdentitiesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedIdentityClient.UserAssignedIdentitiesListByResourceGroup`

```go
ctx := context.TODO()
id := managedidentity.NewResourceGroupID()

// alternatively `client.UserAssignedIdentitiesListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.UserAssignedIdentitiesListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ManagedIdentityClient.UserAssignedIdentitiesListBySubscription`

```go
ctx := context.TODO()
id := managedidentity.NewSubscriptionID()

// alternatively `client.UserAssignedIdentitiesListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.UserAssignedIdentitiesListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ManagedIdentityClient.UserAssignedIdentitiesUpdate`

```go
ctx := context.TODO()
id := managedidentity.NewUserAssignedIdentityID()

payload := managedidentity.IdentityUpdate{
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

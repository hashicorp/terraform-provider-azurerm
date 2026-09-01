
## `github.com/hashicorp/go-azure-sdk/resource-manager/managedidentity/2024-11-30/identities` Documentation

The `identities` SDK allows for interaction with Azure Resource Manager `managedidentity` (API Version `2024-11-30`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/managedidentity/2024-11-30/identities"
```


### Client Initialization

```go
client := identities.NewIdentitiesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `IdentitiesClient.UserAssignedIdentitiesCreateOrUpdate`

```go
ctx := context.TODO()
id := commonids.NewUserAssignedIdentityID("12345678-1234-9876-4563-123456789012", "example-resource-group", "userAssignedIdentityName")

payload := identities.Identity{
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


### Example Usage: `IdentitiesClient.UserAssignedIdentitiesDelete`

```go
ctx := context.TODO()
id := commonids.NewUserAssignedIdentityID("12345678-1234-9876-4563-123456789012", "example-resource-group", "userAssignedIdentityName")

read, err := client.UserAssignedIdentitiesDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IdentitiesClient.UserAssignedIdentitiesGet`

```go
ctx := context.TODO()
id := commonids.NewUserAssignedIdentityID("12345678-1234-9876-4563-123456789012", "example-resource-group", "userAssignedIdentityName")

read, err := client.UserAssignedIdentitiesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IdentitiesClient.UserAssignedIdentitiesListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.UserAssignedIdentitiesListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.UserAssignedIdentitiesListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `IdentitiesClient.UserAssignedIdentitiesListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.UserAssignedIdentitiesListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.UserAssignedIdentitiesListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `IdentitiesClient.UserAssignedIdentitiesUpdate`

```go
ctx := context.TODO()
id := commonids.NewUserAssignedIdentityID("12345678-1234-9876-4563-123456789012", "example-resource-group", "userAssignedIdentityName")

payload := identities.IdentityUpdate{
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

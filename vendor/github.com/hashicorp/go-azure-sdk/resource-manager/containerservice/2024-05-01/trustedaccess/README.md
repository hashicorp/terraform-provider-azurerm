
## `github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2024-05-01/trustedaccess` Documentation

The `trustedaccess` SDK allows for interaction with Azure Resource Manager `containerservice` (API Version `2024-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2024-05-01/trustedaccess"
```


### Client Initialization

```go
client := trustedaccess.NewTrustedAccessClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `TrustedAccessClient.RoleBindingsCreateOrUpdate`

```go
ctx := context.TODO()
id := trustedaccess.NewTrustedAccessRoleBindingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedClusterName", "trustedAccessRoleBindingName")

payload := trustedaccess.TrustedAccessRoleBinding{
	// ...
}


if err := client.RoleBindingsCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `TrustedAccessClient.RoleBindingsDelete`

```go
ctx := context.TODO()
id := trustedaccess.NewTrustedAccessRoleBindingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedClusterName", "trustedAccessRoleBindingName")

if err := client.RoleBindingsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `TrustedAccessClient.RoleBindingsGet`

```go
ctx := context.TODO()
id := trustedaccess.NewTrustedAccessRoleBindingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedClusterName", "trustedAccessRoleBindingName")

read, err := client.RoleBindingsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TrustedAccessClient.RoleBindingsList`

```go
ctx := context.TODO()
id := commonids.NewKubernetesClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedClusterName")

// alternatively `client.RoleBindingsList(ctx, id)` can be used to do batched pagination
items, err := client.RoleBindingsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `TrustedAccessClient.RolesList`

```go
ctx := context.TODO()
id := trustedaccess.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

// alternatively `client.RolesList(ctx, id)` can be used to do batched pagination
items, err := client.RolesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

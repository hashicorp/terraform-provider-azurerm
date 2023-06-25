
## `github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2022-09-02-preview/trustedaccess` Documentation

The `trustedaccess` SDK allows for interaction with the Azure Resource Manager Service `containerservice` (API Version `2022-09-02-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2022-09-02-preview/trustedaccess"
```


### Client Initialization

```go
client := trustedaccess.NewTrustedAccessClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `TrustedAccessClient.RoleBindingsCreateOrUpdate`

```go
ctx := context.TODO()
id := trustedaccess.NewTrustedAccessRoleBindingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedClusterValue", "trustedAccessRoleBindingValue")

payload := trustedaccess.TrustedAccessRoleBinding{
	// ...
}


read, err := client.RoleBindingsCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TrustedAccessClient.RoleBindingsDelete`

```go
ctx := context.TODO()
id := trustedaccess.NewTrustedAccessRoleBindingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedClusterValue", "trustedAccessRoleBindingValue")

read, err := client.RoleBindingsDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TrustedAccessClient.RoleBindingsGet`

```go
ctx := context.TODO()
id := trustedaccess.NewTrustedAccessRoleBindingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedClusterValue", "trustedAccessRoleBindingValue")

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
id := trustedaccess.NewKubernetesClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedClusterValue")

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
id := trustedaccess.NewLocationID("12345678-1234-9876-4563-123456789012", "locationValue")

// alternatively `client.RolesList(ctx, id)` can be used to do batched pagination
items, err := client.RolesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

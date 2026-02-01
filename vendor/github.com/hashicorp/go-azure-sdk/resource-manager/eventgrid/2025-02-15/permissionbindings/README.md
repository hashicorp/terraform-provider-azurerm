
## `github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/permissionbindings` Documentation

The `permissionbindings` SDK allows for interaction with Azure Resource Manager `eventgrid` (API Version `2025-02-15`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/permissionbindings"
```


### Client Initialization

```go
client := permissionbindings.NewPermissionBindingsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PermissionBindingsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := permissionbindings.NewPermissionBindingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "permissionBindingName")

payload := permissionbindings.PermissionBinding{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `PermissionBindingsClient.Delete`

```go
ctx := context.TODO()
id := permissionbindings.NewPermissionBindingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "permissionBindingName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `PermissionBindingsClient.Get`

```go
ctx := context.TODO()
id := permissionbindings.NewPermissionBindingID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "permissionBindingName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PermissionBindingsClient.ListByNamespace`

```go
ctx := context.TODO()
id := permissionbindings.NewNamespaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName")

// alternatively `client.ListByNamespace(ctx, id, permissionbindings.DefaultListByNamespaceOperationOptions())` can be used to do batched pagination
items, err := client.ListByNamespaceComplete(ctx, id, permissionbindings.DefaultListByNamespaceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

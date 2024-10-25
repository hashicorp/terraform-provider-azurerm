
## `github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2022-04-01/roledefinitions` Documentation

The `roledefinitions` SDK allows for interaction with Azure Resource Manager `authorization` (API Version `2022-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2022-04-01/roledefinitions"
```


### Client Initialization

```go
client := roledefinitions.NewRoleDefinitionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `RoleDefinitionsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := roledefinitions.NewScopedRoleDefinitionID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "roleDefinitionId")

payload := roledefinitions.RoleDefinition{
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


### Example Usage: `RoleDefinitionsClient.Delete`

```go
ctx := context.TODO()
id := roledefinitions.NewScopedRoleDefinitionID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "roleDefinitionId")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RoleDefinitionsClient.Get`

```go
ctx := context.TODO()
id := roledefinitions.NewScopedRoleDefinitionID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "roleDefinitionId")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RoleDefinitionsClient.List`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

// alternatively `client.List(ctx, id, roledefinitions.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, roledefinitions.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

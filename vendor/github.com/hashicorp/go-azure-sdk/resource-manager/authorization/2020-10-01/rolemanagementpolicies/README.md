
## `github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2020-10-01/rolemanagementpolicies` Documentation

The `rolemanagementpolicies` SDK allows for interaction with Azure Resource Manager `authorization` (API Version `2020-10-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2020-10-01/rolemanagementpolicies"
```


### Client Initialization

```go
client := rolemanagementpolicies.NewRoleManagementPoliciesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `RoleManagementPoliciesClient.Delete`

```go
ctx := context.TODO()
id := rolemanagementpolicies.NewScopedRoleManagementPolicyID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "roleManagementPolicyName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RoleManagementPoliciesClient.Get`

```go
ctx := context.TODO()
id := rolemanagementpolicies.NewScopedRoleManagementPolicyID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "roleManagementPolicyName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RoleManagementPoliciesClient.ListForScope`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

// alternatively `client.ListForScope(ctx, id, rolemanagementpolicies.DefaultListForScopeOperationOptions())` can be used to do batched pagination
items, err := client.ListForScopeComplete(ctx, id, rolemanagementpolicies.DefaultListForScopeOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `RoleManagementPoliciesClient.Update`

```go
ctx := context.TODO()
id := rolemanagementpolicies.NewScopedRoleManagementPolicyID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "roleManagementPolicyName")

payload := rolemanagementpolicies.RoleManagementPolicy{
	// ...
}


read, err := client.Update(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

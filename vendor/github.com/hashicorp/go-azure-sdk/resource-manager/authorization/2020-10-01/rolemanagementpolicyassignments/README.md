
## `github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2020-10-01/rolemanagementpolicyassignments` Documentation

The `rolemanagementpolicyassignments` SDK allows for interaction with Azure Resource Manager `authorization` (API Version `2020-10-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2020-10-01/rolemanagementpolicyassignments"
```


### Client Initialization

```go
client := rolemanagementpolicyassignments.NewRoleManagementPolicyAssignmentsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `RoleManagementPolicyAssignmentsClient.Create`

```go
ctx := context.TODO()
id := rolemanagementpolicyassignments.NewScopedRoleManagementPolicyAssignmentID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "roleManagementPolicyAssignmentName")

payload := rolemanagementpolicyassignments.RoleManagementPolicyAssignment{
	// ...
}


read, err := client.Create(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RoleManagementPolicyAssignmentsClient.Delete`

```go
ctx := context.TODO()
id := rolemanagementpolicyassignments.NewScopedRoleManagementPolicyAssignmentID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "roleManagementPolicyAssignmentName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RoleManagementPolicyAssignmentsClient.Get`

```go
ctx := context.TODO()
id := rolemanagementpolicyassignments.NewScopedRoleManagementPolicyAssignmentID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "roleManagementPolicyAssignmentName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RoleManagementPolicyAssignmentsClient.ListForScope`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

// alternatively `client.ListForScope(ctx, id)` can be used to do batched pagination
items, err := client.ListForScopeComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

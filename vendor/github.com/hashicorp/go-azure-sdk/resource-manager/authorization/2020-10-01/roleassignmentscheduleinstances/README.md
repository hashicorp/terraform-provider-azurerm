
## `github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2020-10-01/roleassignmentscheduleinstances` Documentation

The `roleassignmentscheduleinstances` SDK allows for interaction with the Azure Resource Manager Service `authorization` (API Version `2020-10-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2020-10-01/roleassignmentscheduleinstances"
```


### Client Initialization

```go
client := roleassignmentscheduleinstances.NewRoleAssignmentScheduleInstancesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `RoleAssignmentScheduleInstancesClient.Get`

```go
ctx := context.TODO()
id := roleassignmentscheduleinstances.NewScopedRoleAssignmentScheduleInstanceID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "roleAssignmentScheduleInstanceValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RoleAssignmentScheduleInstancesClient.ListForScope`

```go
ctx := context.TODO()
id := roleassignmentscheduleinstances.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

// alternatively `client.ListForScope(ctx, id, roleassignmentscheduleinstances.DefaultListForScopeOperationOptions())` can be used to do batched pagination
items, err := client.ListForScopeComplete(ctx, id, roleassignmentscheduleinstances.DefaultListForScopeOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

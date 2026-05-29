
## `github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2020-10-01/roleeligibilityscheduleinstances` Documentation

The `roleeligibilityscheduleinstances` SDK allows for interaction with Azure Resource Manager `authorization` (API Version `2020-10-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2020-10-01/roleeligibilityscheduleinstances"
```


### Client Initialization

```go
client := roleeligibilityscheduleinstances.NewRoleEligibilityScheduleInstancesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `RoleEligibilityScheduleInstancesClient.Get`

```go
ctx := context.TODO()
id := roleeligibilityscheduleinstances.NewScopedRoleEligibilityScheduleInstanceID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "roleEligibilityScheduleInstanceName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RoleEligibilityScheduleInstancesClient.ListForScope`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

// alternatively `client.ListForScope(ctx, id, roleeligibilityscheduleinstances.DefaultListForScopeOperationOptions())` can be used to do batched pagination
items, err := client.ListForScopeComplete(ctx, id, roleeligibilityscheduleinstances.DefaultListForScopeOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


## `github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2020-10-01/roleeligibilityschedulerequests` Documentation

The `roleeligibilityschedulerequests` SDK allows for interaction with the Azure Resource Manager Service `authorization` (API Version `2020-10-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2020-10-01/roleeligibilityschedulerequests"
```


### Client Initialization

```go
client := roleeligibilityschedulerequests.NewRoleEligibilityScheduleRequestsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `RoleEligibilityScheduleRequestsClient.Cancel`

```go
ctx := context.TODO()
id := roleeligibilityschedulerequests.NewScopedRoleEligibilityScheduleRequestID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "roleEligibilityScheduleRequestValue")

read, err := client.Cancel(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RoleEligibilityScheduleRequestsClient.Create`

```go
ctx := context.TODO()
id := roleeligibilityschedulerequests.NewScopedRoleEligibilityScheduleRequestID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "roleEligibilityScheduleRequestValue")

payload := roleeligibilityschedulerequests.RoleEligibilityScheduleRequest{
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


### Example Usage: `RoleEligibilityScheduleRequestsClient.Get`

```go
ctx := context.TODO()
id := roleeligibilityschedulerequests.NewScopedRoleEligibilityScheduleRequestID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "roleEligibilityScheduleRequestValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RoleEligibilityScheduleRequestsClient.ListForScope`

```go
ctx := context.TODO()
id := roleeligibilityschedulerequests.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

// alternatively `client.ListForScope(ctx, id, roleeligibilityschedulerequests.DefaultListForScopeOperationOptions())` can be used to do batched pagination
items, err := client.ListForScopeComplete(ctx, id, roleeligibilityschedulerequests.DefaultListForScopeOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `RoleEligibilityScheduleRequestsClient.Validate`

```go
ctx := context.TODO()
id := roleeligibilityschedulerequests.NewScopedRoleEligibilityScheduleRequestID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "roleEligibilityScheduleRequestValue")

payload := roleeligibilityschedulerequests.RoleEligibilityScheduleRequest{
	// ...
}


read, err := client.Validate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

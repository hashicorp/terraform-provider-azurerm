
## `github.com/hashicorp/go-azure-sdk/resource-manager/managedservices/2019-06-01/registrationassignments` Documentation

The `registrationassignments` SDK allows for interaction with the Azure Resource Manager Service `managedservices` (API Version `2019-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/managedservices/2019-06-01/registrationassignments"
```


### Client Initialization

```go
client := registrationassignments.NewRegistrationAssignmentsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `RegistrationAssignmentsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := registrationassignments.NewScopedRegistrationAssignmentID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "registrationAssignmentIdValue")

payload := registrationassignments.RegistrationAssignment{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `RegistrationAssignmentsClient.Delete`

```go
ctx := context.TODO()
id := registrationassignments.NewScopedRegistrationAssignmentID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "registrationAssignmentIdValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `RegistrationAssignmentsClient.Get`

```go
ctx := context.TODO()
id := registrationassignments.NewScopedRegistrationAssignmentID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "registrationAssignmentIdValue")

read, err := client.Get(ctx, id, registrationassignments.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RegistrationAssignmentsClient.List`

```go
ctx := context.TODO()
id := registrationassignments.NewScopeID()

// alternatively `client.List(ctx, id, registrationassignments.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, registrationassignments.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

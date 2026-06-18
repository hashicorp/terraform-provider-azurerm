
## `github.com/hashicorp/go-azure-sdk/data-plane/synapse/2020-08-01-preview/roleassignments` Documentation

The `roleassignments` SDK allows for interaction with <unknown source data type> `synapse` (API Version `2020-08-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/data-plane/synapse/2020-08-01-preview/roleassignments"
```


### Client Initialization

```go
client := roleassignments.NewRoleAssignmentsClientWithBaseURI("")
client.Client.Authorizer = authorizer
```


### Example Usage: `RoleAssignmentsClient.CreateRoleAssignment`

```go
ctx := context.TODO()
id := roleassignments.NewRoleAssignmentIdID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

payload := roleassignments.RoleAssignmentRequest{
	// ...
}


read, err := client.CreateRoleAssignment(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RoleAssignmentsClient.DeleteRoleAssignmentById`

```go
ctx := context.TODO()
id := roleassignments.NewRoleAssignmentIdID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

read, err := client.DeleteRoleAssignmentById(ctx, id, roleassignments.DefaultDeleteRoleAssignmentByIdOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RoleAssignmentsClient.GetRoleAssignmentById`

```go
ctx := context.TODO()
id := roleassignments.NewRoleAssignmentIdID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

read, err := client.GetRoleAssignmentById(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RoleAssignmentsClient.ListRoleAssignments`

```go
ctx := context.TODO()


read, err := client.ListRoleAssignments(ctx, roleassignments.DefaultListRoleAssignmentsOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

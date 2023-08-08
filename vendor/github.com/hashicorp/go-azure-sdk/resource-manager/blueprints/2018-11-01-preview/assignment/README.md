
## `github.com/hashicorp/go-azure-sdk/resource-manager/blueprints/2018-11-01-preview/assignment` Documentation

The `assignment` SDK allows for interaction with the Azure Resource Manager Service `blueprints` (API Version `2018-11-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/blueprints/2018-11-01-preview/assignment"
```


### Client Initialization

```go
client := assignment.NewAssignmentClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AssignmentClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := assignment.NewScopedBlueprintAssignmentID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "blueprintAssignmentValue")

payload := assignment.Assignment{
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


### Example Usage: `AssignmentClient.Delete`

```go
ctx := context.TODO()
id := assignment.NewScopedBlueprintAssignmentID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "blueprintAssignmentValue")

read, err := client.Delete(ctx, id, assignment.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AssignmentClient.Get`

```go
ctx := context.TODO()
id := assignment.NewScopedBlueprintAssignmentID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "blueprintAssignmentValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AssignmentClient.List`

```go
ctx := context.TODO()
id := assignment.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


## `github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-06-01/policyassignments` Documentation

The `policyassignments` SDK allows for interaction with Azure Resource Manager `resources` (API Version `2022-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-06-01/policyassignments"
```


### Client Initialization

```go
client := policyassignments.NewPolicyAssignmentsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PolicyAssignmentsClient.Create`

```go
ctx := context.TODO()
id := policyassignments.NewScopedPolicyAssignmentID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "policyAssignmentName")

payload := policyassignments.PolicyAssignment{
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


### Example Usage: `PolicyAssignmentsClient.CreateById`

```go
ctx := context.TODO()
id := policyassignments.NewPolicyAssignmentIdID("policyAssignmentId")

payload := policyassignments.PolicyAssignment{
	// ...
}


read, err := client.CreateById(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PolicyAssignmentsClient.Delete`

```go
ctx := context.TODO()
id := policyassignments.NewScopedPolicyAssignmentID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "policyAssignmentName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PolicyAssignmentsClient.DeleteById`

```go
ctx := context.TODO()
id := policyassignments.NewPolicyAssignmentIdID("policyAssignmentId")

read, err := client.DeleteById(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PolicyAssignmentsClient.Get`

```go
ctx := context.TODO()
id := policyassignments.NewScopedPolicyAssignmentID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "policyAssignmentName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PolicyAssignmentsClient.GetById`

```go
ctx := context.TODO()
id := policyassignments.NewPolicyAssignmentIdID("policyAssignmentId")

read, err := client.GetById(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PolicyAssignmentsClient.List`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.List(ctx, id, policyassignments.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, policyassignments.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PolicyAssignmentsClient.ListForManagementGroup`

```go
ctx := context.TODO()
id := commonids.NewManagementGroupID("groupId")

// alternatively `client.ListForManagementGroup(ctx, id, policyassignments.DefaultListForManagementGroupOperationOptions())` can be used to do batched pagination
items, err := client.ListForManagementGroupComplete(ctx, id, policyassignments.DefaultListForManagementGroupOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PolicyAssignmentsClient.ListForResource`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

// alternatively `client.ListForResource(ctx, id, policyassignments.DefaultListForResourceOperationOptions())` can be used to do batched pagination
items, err := client.ListForResourceComplete(ctx, id, policyassignments.DefaultListForResourceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PolicyAssignmentsClient.ListForResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListForResourceGroup(ctx, id, policyassignments.DefaultListForResourceGroupOperationOptions())` can be used to do batched pagination
items, err := client.ListForResourceGroupComplete(ctx, id, policyassignments.DefaultListForResourceGroupOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PolicyAssignmentsClient.Update`

```go
ctx := context.TODO()
id := policyassignments.NewScopedPolicyAssignmentID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "policyAssignmentName")

payload := policyassignments.PolicyAssignmentUpdate{
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


### Example Usage: `PolicyAssignmentsClient.UpdateById`

```go
ctx := context.TODO()
id := policyassignments.NewPolicyAssignmentIdID("policyAssignmentId")

payload := policyassignments.PolicyAssignmentUpdate{
	// ...
}


read, err := client.UpdateById(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

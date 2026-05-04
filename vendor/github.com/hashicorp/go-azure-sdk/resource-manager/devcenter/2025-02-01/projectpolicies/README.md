
## `github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/projectpolicies` Documentation

The `projectpolicies` SDK allows for interaction with Azure Resource Manager `devcenter` (API Version `2025-02-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/projectpolicies"
```


### Client Initialization

```go
client := projectpolicies.NewProjectPoliciesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ProjectPoliciesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := projectpolicies.NewProjectPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "devCenterName", "projectPolicyName")

payload := projectpolicies.ProjectPolicy{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ProjectPoliciesClient.Delete`

```go
ctx := context.TODO()
id := projectpolicies.NewProjectPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "devCenterName", "projectPolicyName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ProjectPoliciesClient.Get`

```go
ctx := context.TODO()
id := projectpolicies.NewProjectPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "devCenterName", "projectPolicyName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProjectPoliciesClient.ListByDevCenter`

```go
ctx := context.TODO()
id := projectpolicies.NewDevCenterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "devCenterName")

// alternatively `client.ListByDevCenter(ctx, id, projectpolicies.DefaultListByDevCenterOperationOptions())` can be used to do batched pagination
items, err := client.ListByDevCenterComplete(ctx, id, projectpolicies.DefaultListByDevCenterOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ProjectPoliciesClient.Update`

```go
ctx := context.TODO()
id := projectpolicies.NewProjectPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "devCenterName", "projectPolicyName")

payload := projectpolicies.ProjectPolicyUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

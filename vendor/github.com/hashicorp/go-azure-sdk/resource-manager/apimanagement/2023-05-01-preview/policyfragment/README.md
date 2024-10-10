
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2023-05-01-preview/policyfragment` Documentation

The `policyfragment` SDK allows for interaction with the Azure Resource Manager Service `apimanagement` (API Version `2023-05-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2023-05-01-preview/policyfragment"
```


### Client Initialization

```go
client := policyfragment.NewPolicyFragmentClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PolicyFragmentClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := policyfragment.NewPolicyFragmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "policyFragmentValue")

payload := policyfragment.PolicyFragmentContract{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload, policyfragment.DefaultCreateOrUpdateOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `PolicyFragmentClient.Delete`

```go
ctx := context.TODO()
id := policyfragment.NewPolicyFragmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "policyFragmentValue")

read, err := client.Delete(ctx, id, policyfragment.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PolicyFragmentClient.Get`

```go
ctx := context.TODO()
id := policyfragment.NewPolicyFragmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "policyFragmentValue")

read, err := client.Get(ctx, id, policyfragment.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PolicyFragmentClient.GetEntityTag`

```go
ctx := context.TODO()
id := policyfragment.NewPolicyFragmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "policyFragmentValue")

read, err := client.GetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PolicyFragmentClient.ListByService`

```go
ctx := context.TODO()
id := policyfragment.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue")

// alternatively `client.ListByService(ctx, id, policyfragment.DefaultListByServiceOperationOptions())` can be used to do batched pagination
items, err := client.ListByServiceComplete(ctx, id, policyfragment.DefaultListByServiceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PolicyFragmentClient.ListReferences`

```go
ctx := context.TODO()
id := policyfragment.NewPolicyFragmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "policyFragmentValue")

// alternatively `client.ListReferences(ctx, id, policyfragment.DefaultListReferencesOperationOptions())` can be used to do batched pagination
items, err := client.ListReferencesComplete(ctx, id, policyfragment.DefaultListReferencesOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PolicyFragmentClient.WorkspacePolicyFragmentCreateOrUpdate`

```go
ctx := context.TODO()
id := policyfragment.NewWorkspacePolicyFragmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "policyFragmentValue")

payload := policyfragment.PolicyFragmentContract{
	// ...
}


if err := client.WorkspacePolicyFragmentCreateOrUpdateThenPoll(ctx, id, payload, policyfragment.DefaultWorkspacePolicyFragmentCreateOrUpdateOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `PolicyFragmentClient.WorkspacePolicyFragmentDelete`

```go
ctx := context.TODO()
id := policyfragment.NewWorkspacePolicyFragmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "policyFragmentValue")

read, err := client.WorkspacePolicyFragmentDelete(ctx, id, policyfragment.DefaultWorkspacePolicyFragmentDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PolicyFragmentClient.WorkspacePolicyFragmentGet`

```go
ctx := context.TODO()
id := policyfragment.NewWorkspacePolicyFragmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "policyFragmentValue")

read, err := client.WorkspacePolicyFragmentGet(ctx, id, policyfragment.DefaultWorkspacePolicyFragmentGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PolicyFragmentClient.WorkspacePolicyFragmentGetEntityTag`

```go
ctx := context.TODO()
id := policyfragment.NewWorkspacePolicyFragmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "policyFragmentValue")

read, err := client.WorkspacePolicyFragmentGetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PolicyFragmentClient.WorkspacePolicyFragmentListByService`

```go
ctx := context.TODO()
id := policyfragment.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue")

// alternatively `client.WorkspacePolicyFragmentListByService(ctx, id, policyfragment.DefaultWorkspacePolicyFragmentListByServiceOperationOptions())` can be used to do batched pagination
items, err := client.WorkspacePolicyFragmentListByServiceComplete(ctx, id, policyfragment.DefaultWorkspacePolicyFragmentListByServiceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PolicyFragmentClient.WorkspacePolicyFragmentListReferences`

```go
ctx := context.TODO()
id := policyfragment.NewWorkspacePolicyFragmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "policyFragmentValue")

// alternatively `client.WorkspacePolicyFragmentListReferences(ctx, id, policyfragment.DefaultWorkspacePolicyFragmentListReferencesOperationOptions())` can be used to do batched pagination
items, err := client.WorkspacePolicyFragmentListReferencesComplete(ctx, id, policyfragment.DefaultWorkspacePolicyFragmentListReferencesOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2023-05-01-preview/apipolicy` Documentation

The `apipolicy` SDK allows for interaction with the Azure Resource Manager Service `apimanagement` (API Version `2023-05-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2023-05-01-preview/apipolicy"
```


### Client Initialization

```go
client := apipolicy.NewApiPolicyClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ApiPolicyClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := apipolicy.NewApiID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "apiIdValue")

payload := apipolicy.PolicyContract{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, apipolicy.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiPolicyClient.Delete`

```go
ctx := context.TODO()
id := apipolicy.NewApiID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "apiIdValue")

read, err := client.Delete(ctx, id, apipolicy.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiPolicyClient.Get`

```go
ctx := context.TODO()
id := apipolicy.NewApiID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "apiIdValue")

read, err := client.Get(ctx, id, apipolicy.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiPolicyClient.GetEntityTag`

```go
ctx := context.TODO()
id := apipolicy.NewApiID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "apiIdValue")

read, err := client.GetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiPolicyClient.ListByApi`

```go
ctx := context.TODO()
id := apipolicy.NewApiID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "apiIdValue")

// alternatively `client.ListByApi(ctx, id)` can be used to do batched pagination
items, err := client.ListByApiComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ApiPolicyClient.WorkspaceApiPolicyCreateOrUpdate`

```go
ctx := context.TODO()
id := apipolicy.NewWorkspaceApiID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "apiIdValue")

payload := apipolicy.PolicyContract{
	// ...
}


read, err := client.WorkspaceApiPolicyCreateOrUpdate(ctx, id, payload, apipolicy.DefaultWorkspaceApiPolicyCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiPolicyClient.WorkspaceApiPolicyDelete`

```go
ctx := context.TODO()
id := apipolicy.NewWorkspaceApiID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "apiIdValue")

read, err := client.WorkspaceApiPolicyDelete(ctx, id, apipolicy.DefaultWorkspaceApiPolicyDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiPolicyClient.WorkspaceApiPolicyGet`

```go
ctx := context.TODO()
id := apipolicy.NewWorkspaceApiID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "apiIdValue")

read, err := client.WorkspaceApiPolicyGet(ctx, id, apipolicy.DefaultWorkspaceApiPolicyGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiPolicyClient.WorkspaceApiPolicyGetEntityTag`

```go
ctx := context.TODO()
id := apipolicy.NewWorkspaceApiID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "apiIdValue")

read, err := client.WorkspaceApiPolicyGetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiPolicyClient.WorkspaceApiPolicyListByApi`

```go
ctx := context.TODO()
id := apipolicy.NewWorkspaceApiID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "apiIdValue")

// alternatively `client.WorkspaceApiPolicyListByApi(ctx, id)` can be used to do batched pagination
items, err := client.WorkspaceApiPolicyListByApiComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

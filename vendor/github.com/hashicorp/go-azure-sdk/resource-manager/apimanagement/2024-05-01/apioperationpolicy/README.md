
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apioperationpolicy` Documentation

The `apioperationpolicy` SDK allows for interaction with Azure Resource Manager `apimanagement` (API Version `2024-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apioperationpolicy"
```


### Client Initialization

```go
client := apioperationpolicy.NewApiOperationPolicyClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ApiOperationPolicyClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := apioperationpolicy.NewOperationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId", "operationId")

payload := apioperationpolicy.PolicyContract{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, apioperationpolicy.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiOperationPolicyClient.Delete`

```go
ctx := context.TODO()
id := apioperationpolicy.NewOperationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId", "operationId")

read, err := client.Delete(ctx, id, apioperationpolicy.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiOperationPolicyClient.Get`

```go
ctx := context.TODO()
id := apioperationpolicy.NewOperationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId", "operationId")

read, err := client.Get(ctx, id, apioperationpolicy.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiOperationPolicyClient.GetEntityTag`

```go
ctx := context.TODO()
id := apioperationpolicy.NewOperationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId", "operationId")

read, err := client.GetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiOperationPolicyClient.ListByOperation`

```go
ctx := context.TODO()
id := apioperationpolicy.NewOperationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId", "operationId")

// alternatively `client.ListByOperation(ctx, id)` can be used to do batched pagination
items, err := client.ListByOperationComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ApiOperationPolicyClient.WorkspaceApiOperationPolicyCreateOrUpdate`

```go
ctx := context.TODO()
id := apioperationpolicy.NewApiOperationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId", "apiId", "operationId")

payload := apioperationpolicy.PolicyContract{
	// ...
}


read, err := client.WorkspaceApiOperationPolicyCreateOrUpdate(ctx, id, payload, apioperationpolicy.DefaultWorkspaceApiOperationPolicyCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiOperationPolicyClient.WorkspaceApiOperationPolicyDelete`

```go
ctx := context.TODO()
id := apioperationpolicy.NewApiOperationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId", "apiId", "operationId")

read, err := client.WorkspaceApiOperationPolicyDelete(ctx, id, apioperationpolicy.DefaultWorkspaceApiOperationPolicyDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiOperationPolicyClient.WorkspaceApiOperationPolicyGet`

```go
ctx := context.TODO()
id := apioperationpolicy.NewApiOperationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId", "apiId", "operationId")

read, err := client.WorkspaceApiOperationPolicyGet(ctx, id, apioperationpolicy.DefaultWorkspaceApiOperationPolicyGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiOperationPolicyClient.WorkspaceApiOperationPolicyGetEntityTag`

```go
ctx := context.TODO()
id := apioperationpolicy.NewApiOperationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId", "apiId", "operationId")

read, err := client.WorkspaceApiOperationPolicyGetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiOperationPolicyClient.WorkspaceApiOperationPolicyListByOperation`

```go
ctx := context.TODO()
id := apioperationpolicy.NewApiOperationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId", "apiId", "operationId")

// alternatively `client.WorkspaceApiOperationPolicyListByOperation(ctx, id)` can be used to do batched pagination
items, err := client.WorkspaceApiOperationPolicyListByOperationComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

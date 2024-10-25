
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2023-05-01-preview/apioperation` Documentation

The `apioperation` SDK allows for interaction with the Azure Resource Manager Service `apimanagement` (API Version `2023-05-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2023-05-01-preview/apioperation"
```


### Client Initialization

```go
client := apioperation.NewApiOperationClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ApiOperationClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := apioperation.NewOperationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "apiIdValue", "operationIdValue")

payload := apioperation.OperationContract{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, apioperation.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiOperationClient.Delete`

```go
ctx := context.TODO()
id := apioperation.NewOperationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "apiIdValue", "operationIdValue")

read, err := client.Delete(ctx, id, apioperation.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiOperationClient.Get`

```go
ctx := context.TODO()
id := apioperation.NewOperationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "apiIdValue", "operationIdValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiOperationClient.GetEntityTag`

```go
ctx := context.TODO()
id := apioperation.NewOperationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "apiIdValue", "operationIdValue")

read, err := client.GetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiOperationClient.ListByApi`

```go
ctx := context.TODO()
id := apioperation.NewApiID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "apiIdValue")

// alternatively `client.ListByApi(ctx, id, apioperation.DefaultListByApiOperationOptions())` can be used to do batched pagination
items, err := client.ListByApiComplete(ctx, id, apioperation.DefaultListByApiOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ApiOperationClient.Update`

```go
ctx := context.TODO()
id := apioperation.NewOperationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "apiIdValue", "operationIdValue")

payload := apioperation.OperationUpdateContract{
	// ...
}


read, err := client.Update(ctx, id, payload, apioperation.DefaultUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiOperationClient.WorkspaceApiOperationCreateOrUpdate`

```go
ctx := context.TODO()
id := apioperation.NewApiOperationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "apiIdValue", "operationIdValue")

payload := apioperation.OperationContract{
	// ...
}


read, err := client.WorkspaceApiOperationCreateOrUpdate(ctx, id, payload, apioperation.DefaultWorkspaceApiOperationCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiOperationClient.WorkspaceApiOperationDelete`

```go
ctx := context.TODO()
id := apioperation.NewApiOperationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "apiIdValue", "operationIdValue")

read, err := client.WorkspaceApiOperationDelete(ctx, id, apioperation.DefaultWorkspaceApiOperationDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiOperationClient.WorkspaceApiOperationGet`

```go
ctx := context.TODO()
id := apioperation.NewApiOperationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "apiIdValue", "operationIdValue")

read, err := client.WorkspaceApiOperationGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiOperationClient.WorkspaceApiOperationGetEntityTag`

```go
ctx := context.TODO()
id := apioperation.NewApiOperationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "apiIdValue", "operationIdValue")

read, err := client.WorkspaceApiOperationGetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiOperationClient.WorkspaceApiOperationListByApi`

```go
ctx := context.TODO()
id := apioperation.NewWorkspaceApiID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "apiIdValue")

// alternatively `client.WorkspaceApiOperationListByApi(ctx, id, apioperation.DefaultWorkspaceApiOperationListByApiOperationOptions())` can be used to do batched pagination
items, err := client.WorkspaceApiOperationListByApiComplete(ctx, id, apioperation.DefaultWorkspaceApiOperationListByApiOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ApiOperationClient.WorkspaceApiOperationUpdate`

```go
ctx := context.TODO()
id := apioperation.NewApiOperationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "apiIdValue", "operationIdValue")

payload := apioperation.OperationUpdateContract{
	// ...
}


read, err := client.WorkspaceApiOperationUpdate(ctx, id, payload, apioperation.DefaultWorkspaceApiOperationUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

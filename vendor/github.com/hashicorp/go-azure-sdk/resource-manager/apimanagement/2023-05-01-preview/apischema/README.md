
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2023-05-01-preview/apischema` Documentation

The `apischema` SDK allows for interaction with the Azure Resource Manager Service `apimanagement` (API Version `2023-05-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2023-05-01-preview/apischema"
```


### Client Initialization

```go
client := apischema.NewApiSchemaClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ApiSchemaClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := apischema.NewApiSchemaID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "apiIdValue", "schemaIdValue")

payload := apischema.SchemaContract{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload, apischema.DefaultCreateOrUpdateOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `ApiSchemaClient.Delete`

```go
ctx := context.TODO()
id := apischema.NewApiSchemaID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "apiIdValue", "schemaIdValue")

read, err := client.Delete(ctx, id, apischema.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiSchemaClient.Get`

```go
ctx := context.TODO()
id := apischema.NewApiSchemaID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "apiIdValue", "schemaIdValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiSchemaClient.GetEntityTag`

```go
ctx := context.TODO()
id := apischema.NewApiSchemaID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "apiIdValue", "schemaIdValue")

read, err := client.GetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiSchemaClient.ListByApi`

```go
ctx := context.TODO()
id := apischema.NewApiID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "apiIdValue")

// alternatively `client.ListByApi(ctx, id, apischema.DefaultListByApiOperationOptions())` can be used to do batched pagination
items, err := client.ListByApiComplete(ctx, id, apischema.DefaultListByApiOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ApiSchemaClient.WorkspaceApiSchemaCreateOrUpdate`

```go
ctx := context.TODO()
id := apischema.NewWorkspaceApiSchemaID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "apiIdValue", "schemaIdValue")

payload := apischema.SchemaContract{
	// ...
}


if err := client.WorkspaceApiSchemaCreateOrUpdateThenPoll(ctx, id, payload, apischema.DefaultWorkspaceApiSchemaCreateOrUpdateOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `ApiSchemaClient.WorkspaceApiSchemaDelete`

```go
ctx := context.TODO()
id := apischema.NewWorkspaceApiSchemaID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "apiIdValue", "schemaIdValue")

read, err := client.WorkspaceApiSchemaDelete(ctx, id, apischema.DefaultWorkspaceApiSchemaDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiSchemaClient.WorkspaceApiSchemaGet`

```go
ctx := context.TODO()
id := apischema.NewWorkspaceApiSchemaID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "apiIdValue", "schemaIdValue")

read, err := client.WorkspaceApiSchemaGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiSchemaClient.WorkspaceApiSchemaGetEntityTag`

```go
ctx := context.TODO()
id := apischema.NewWorkspaceApiSchemaID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "apiIdValue", "schemaIdValue")

read, err := client.WorkspaceApiSchemaGetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiSchemaClient.WorkspaceApiSchemaListByApi`

```go
ctx := context.TODO()
id := apischema.NewWorkspaceApiID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "apiIdValue")

// alternatively `client.WorkspaceApiSchemaListByApi(ctx, id, apischema.DefaultWorkspaceApiSchemaListByApiOperationOptions())` can be used to do batched pagination
items, err := client.WorkspaceApiSchemaListByApiComplete(ctx, id, apischema.DefaultWorkspaceApiSchemaListByApiOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2023-05-01-preview/schema` Documentation

The `schema` SDK allows for interaction with the Azure Resource Manager Service `apimanagement` (API Version `2023-05-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2023-05-01-preview/schema"
```


### Client Initialization

```go
client := schema.NewSchemaClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SchemaClient.GlobalSchemaCreateOrUpdate`

```go
ctx := context.TODO()
id := schema.NewSchemaID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "schemaIdValue")

payload := schema.GlobalSchemaContract{
	// ...
}


if err := client.GlobalSchemaCreateOrUpdateThenPoll(ctx, id, payload, schema.DefaultGlobalSchemaCreateOrUpdateOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `SchemaClient.GlobalSchemaDelete`

```go
ctx := context.TODO()
id := schema.NewSchemaID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "schemaIdValue")

read, err := client.GlobalSchemaDelete(ctx, id, schema.DefaultGlobalSchemaDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SchemaClient.GlobalSchemaGet`

```go
ctx := context.TODO()
id := schema.NewSchemaID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "schemaIdValue")

read, err := client.GlobalSchemaGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SchemaClient.GlobalSchemaGetEntityTag`

```go
ctx := context.TODO()
id := schema.NewSchemaID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "schemaIdValue")

read, err := client.GlobalSchemaGetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SchemaClient.GlobalSchemaListByService`

```go
ctx := context.TODO()
id := schema.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue")

// alternatively `client.GlobalSchemaListByService(ctx, id, schema.DefaultGlobalSchemaListByServiceOperationOptions())` can be used to do batched pagination
items, err := client.GlobalSchemaListByServiceComplete(ctx, id, schema.DefaultGlobalSchemaListByServiceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `SchemaClient.WorkspaceGlobalSchemaCreateOrUpdate`

```go
ctx := context.TODO()
id := schema.NewWorkspaceSchemaID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "schemaIdValue")

payload := schema.GlobalSchemaContract{
	// ...
}


if err := client.WorkspaceGlobalSchemaCreateOrUpdateThenPoll(ctx, id, payload, schema.DefaultWorkspaceGlobalSchemaCreateOrUpdateOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `SchemaClient.WorkspaceGlobalSchemaDelete`

```go
ctx := context.TODO()
id := schema.NewWorkspaceSchemaID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "schemaIdValue")

read, err := client.WorkspaceGlobalSchemaDelete(ctx, id, schema.DefaultWorkspaceGlobalSchemaDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SchemaClient.WorkspaceGlobalSchemaGet`

```go
ctx := context.TODO()
id := schema.NewWorkspaceSchemaID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "schemaIdValue")

read, err := client.WorkspaceGlobalSchemaGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SchemaClient.WorkspaceGlobalSchemaGetEntityTag`

```go
ctx := context.TODO()
id := schema.NewWorkspaceSchemaID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "schemaIdValue")

read, err := client.WorkspaceGlobalSchemaGetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SchemaClient.WorkspaceGlobalSchemaListByService`

```go
ctx := context.TODO()
id := schema.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue")

// alternatively `client.WorkspaceGlobalSchemaListByService(ctx, id, schema.DefaultWorkspaceGlobalSchemaListByServiceOperationOptions())` can be used to do batched pagination
items, err := client.WorkspaceGlobalSchemaListByServiceComplete(ctx, id, schema.DefaultWorkspaceGlobalSchemaListByServiceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

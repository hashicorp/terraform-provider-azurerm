
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/schema` Documentation

The `schema` SDK allows for interaction with the Azure Resource Manager Service `apimanagement` (API Version `2021-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/schema"
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

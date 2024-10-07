
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apischema` Documentation

The `apischema` SDK allows for interaction with Azure Resource Manager `apimanagement` (API Version `2022-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apischema"
```


### Client Initialization

```go
client := apischema.NewApiSchemaClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ApiSchemaClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := apischema.NewApiSchemaID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId", "schemaId")

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
id := apischema.NewApiSchemaID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId", "schemaId")

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
id := apischema.NewApiSchemaID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId", "schemaId")

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
id := apischema.NewApiSchemaID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId", "schemaId")

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
id := apischema.NewApiID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId")

// alternatively `client.ListByApi(ctx, id, apischema.DefaultListByApiOperationOptions())` can be used to do batched pagination
items, err := client.ListByApiComplete(ctx, id, apischema.DefaultListByApiOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

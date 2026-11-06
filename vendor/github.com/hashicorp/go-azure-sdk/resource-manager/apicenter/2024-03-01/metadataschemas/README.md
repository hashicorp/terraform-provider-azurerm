
## `github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/metadataschemas` Documentation

The `metadataschemas` SDK allows for interaction with Azure Resource Manager `apicenter` (API Version `2024-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/metadataschemas"
```


### Client Initialization

```go
client := metadataschemas.NewMetadataSchemasClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `MetadataSchemasClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := metadataschemas.NewMetadataSchemaID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "metadataSchemaName")

payload := metadataschemas.MetadataSchema{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `MetadataSchemasClient.Delete`

```go
ctx := context.TODO()
id := metadataschemas.NewMetadataSchemaID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "metadataSchemaName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `MetadataSchemasClient.Get`

```go
ctx := context.TODO()
id := metadataschemas.NewMetadataSchemaID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "metadataSchemaName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `MetadataSchemasClient.Head`

```go
ctx := context.TODO()
id := metadataschemas.NewMetadataSchemaID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "metadataSchemaName")

read, err := client.Head(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `MetadataSchemasClient.List`

```go
ctx := context.TODO()
id := metadataschemas.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName")

// alternatively `client.List(ctx, id, metadataschemas.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, metadataschemas.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

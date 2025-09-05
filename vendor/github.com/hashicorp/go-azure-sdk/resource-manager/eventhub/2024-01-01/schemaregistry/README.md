
## `github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/schemaregistry` Documentation

The `schemaregistry` SDK allows for interaction with Azure Resource Manager `eventhub` (API Version `2024-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/schemaregistry"
```


### Client Initialization

```go
client := schemaregistry.NewSchemaRegistryClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SchemaRegistryClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := schemaregistry.NewSchemaGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "schemaGroupName")

payload := schemaregistry.SchemaGroup{
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


### Example Usage: `SchemaRegistryClient.Delete`

```go
ctx := context.TODO()
id := schemaregistry.NewSchemaGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "schemaGroupName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SchemaRegistryClient.Get`

```go
ctx := context.TODO()
id := schemaregistry.NewSchemaGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "schemaGroupName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SchemaRegistryClient.ListByNamespace`

```go
ctx := context.TODO()
id := schemaregistry.NewNamespaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName")

// alternatively `client.ListByNamespace(ctx, id, schemaregistry.DefaultListByNamespaceOperationOptions())` can be used to do batched pagination
items, err := client.ListByNamespaceComplete(ctx, id, schemaregistry.DefaultListByNamespaceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

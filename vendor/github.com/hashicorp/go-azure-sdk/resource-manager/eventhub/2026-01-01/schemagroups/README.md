
## `github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2026-01-01/schemagroups` Documentation

The `schemagroups` SDK allows for interaction with Azure Resource Manager `eventhub` (API Version `2026-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2026-01-01/schemagroups"
```


### Client Initialization

```go
client := schemagroups.NewSchemaGroupsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SchemaGroupsClient.SchemaRegistryCreateOrUpdate`

```go
ctx := context.TODO()
id := schemagroups.NewSchemaGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "schemaGroupName")

payload := schemagroups.SchemaGroup{
	// ...
}


read, err := client.SchemaRegistryCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SchemaGroupsClient.SchemaRegistryDelete`

```go
ctx := context.TODO()
id := schemagroups.NewSchemaGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "schemaGroupName")

read, err := client.SchemaRegistryDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SchemaGroupsClient.SchemaRegistryGet`

```go
ctx := context.TODO()
id := schemagroups.NewSchemaGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "schemaGroupName")

read, err := client.SchemaRegistryGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SchemaGroupsClient.SchemaRegistryListByNamespace`

```go
ctx := context.TODO()
id := schemagroups.NewNamespaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName")

// alternatively `client.SchemaRegistryListByNamespace(ctx, id, schemagroups.DefaultSchemaRegistryListByNamespaceOperationOptions())` can be used to do batched pagination
items, err := client.SchemaRegistryListByNamespaceComplete(ctx, id, schemagroups.DefaultSchemaRegistryListByNamespaceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


## `github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccountschemas` Documentation

The `integrationaccountschemas` SDK allows for interaction with the Azure Resource Manager Service `logic` (API Version `2019-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccountschemas"
```


### Client Initialization

```go
client := integrationaccountschemas.NewIntegrationAccountSchemasClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `IntegrationAccountSchemasClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := integrationaccountschemas.NewSchemaID("12345678-1234-9876-4563-123456789012", "example-resource-group", "integrationAccountValue", "schemaValue")

payload := integrationaccountschemas.IntegrationAccountSchema{
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


### Example Usage: `IntegrationAccountSchemasClient.Delete`

```go
ctx := context.TODO()
id := integrationaccountschemas.NewSchemaID("12345678-1234-9876-4563-123456789012", "example-resource-group", "integrationAccountValue", "schemaValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IntegrationAccountSchemasClient.Get`

```go
ctx := context.TODO()
id := integrationaccountschemas.NewSchemaID("12345678-1234-9876-4563-123456789012", "example-resource-group", "integrationAccountValue", "schemaValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IntegrationAccountSchemasClient.List`

```go
ctx := context.TODO()
id := integrationaccountschemas.NewIntegrationAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "integrationAccountValue")

// alternatively `client.List(ctx, id, integrationaccountschemas.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, integrationaccountschemas.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `IntegrationAccountSchemasClient.ListContentCallbackUrl`

```go
ctx := context.TODO()
id := integrationaccountschemas.NewSchemaID("12345678-1234-9876-4563-123456789012", "example-resource-group", "integrationAccountValue", "schemaValue")

payload := integrationaccountschemas.GetCallbackUrlParameters{
	// ...
}


read, err := client.ListContentCallbackUrl(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

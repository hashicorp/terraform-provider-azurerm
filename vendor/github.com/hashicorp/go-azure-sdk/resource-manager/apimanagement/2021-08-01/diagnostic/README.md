
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/diagnostic` Documentation

The `diagnostic` SDK allows for interaction with the Azure Resource Manager Service `apimanagement` (API Version `2021-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/diagnostic"
```


### Client Initialization

```go
client := diagnostic.NewDiagnosticClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DiagnosticClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := diagnostic.NewDiagnosticID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "diagnosticIdValue")

payload := diagnostic.DiagnosticContract{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, diagnostic.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DiagnosticClient.Delete`

```go
ctx := context.TODO()
id := diagnostic.NewDiagnosticID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "diagnosticIdValue")

read, err := client.Delete(ctx, id, diagnostic.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DiagnosticClient.Get`

```go
ctx := context.TODO()
id := diagnostic.NewDiagnosticID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "diagnosticIdValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DiagnosticClient.GetEntityTag`

```go
ctx := context.TODO()
id := diagnostic.NewDiagnosticID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "diagnosticIdValue")

read, err := client.GetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DiagnosticClient.ListByService`

```go
ctx := context.TODO()
id := diagnostic.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue")

// alternatively `client.ListByService(ctx, id, diagnostic.DefaultListByServiceOperationOptions())` can be used to do batched pagination
items, err := client.ListByServiceComplete(ctx, id, diagnostic.DefaultListByServiceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DiagnosticClient.Update`

```go
ctx := context.TODO()
id := diagnostic.NewDiagnosticID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "diagnosticIdValue")

payload := diagnostic.DiagnosticContract{
	// ...
}


read, err := client.Update(ctx, id, payload, diagnostic.DefaultUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

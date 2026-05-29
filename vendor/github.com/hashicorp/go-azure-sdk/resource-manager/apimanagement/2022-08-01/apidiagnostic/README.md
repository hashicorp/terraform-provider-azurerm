
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apidiagnostic` Documentation

The `apidiagnostic` SDK allows for interaction with Azure Resource Manager `apimanagement` (API Version `2022-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apidiagnostic"
```


### Client Initialization

```go
client := apidiagnostic.NewApiDiagnosticClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ApiDiagnosticClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := apidiagnostic.NewApiDiagnosticID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId", "diagnosticId")

payload := apidiagnostic.DiagnosticContract{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, apidiagnostic.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiDiagnosticClient.Delete`

```go
ctx := context.TODO()
id := apidiagnostic.NewApiDiagnosticID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId", "diagnosticId")

read, err := client.Delete(ctx, id, apidiagnostic.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiDiagnosticClient.Get`

```go
ctx := context.TODO()
id := apidiagnostic.NewApiDiagnosticID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId", "diagnosticId")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiDiagnosticClient.GetEntityTag`

```go
ctx := context.TODO()
id := apidiagnostic.NewApiDiagnosticID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId", "diagnosticId")

read, err := client.GetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiDiagnosticClient.ListByService`

```go
ctx := context.TODO()
id := apidiagnostic.NewApiID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId")

// alternatively `client.ListByService(ctx, id, apidiagnostic.DefaultListByServiceOperationOptions())` can be used to do batched pagination
items, err := client.ListByServiceComplete(ctx, id, apidiagnostic.DefaultListByServiceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ApiDiagnosticClient.Update`

```go
ctx := context.TODO()
id := apidiagnostic.NewApiDiagnosticID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId", "diagnosticId")

payload := apidiagnostic.DiagnosticContract{
	// ...
}


read, err := client.Update(ctx, id, payload, apidiagnostic.DefaultUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

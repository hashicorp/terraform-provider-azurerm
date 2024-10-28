
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/diagnostic` Documentation

The `diagnostic` SDK allows for interaction with Azure Resource Manager `apimanagement` (API Version `2024-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/diagnostic"
```


### Client Initialization

```go
client := diagnostic.NewDiagnosticClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DiagnosticClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := diagnostic.NewDiagnosticID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "diagnosticId")

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
id := diagnostic.NewDiagnosticID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "diagnosticId")

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
id := diagnostic.NewDiagnosticID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "diagnosticId")

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
id := diagnostic.NewDiagnosticID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "diagnosticId")

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
id := diagnostic.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName")

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
id := diagnostic.NewDiagnosticID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "diagnosticId")

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


### Example Usage: `DiagnosticClient.WorkspaceDiagnosticCreateOrUpdate`

```go
ctx := context.TODO()
id := diagnostic.NewWorkspaceDiagnosticID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId", "diagnosticId")

payload := diagnostic.DiagnosticContract{
	// ...
}


read, err := client.WorkspaceDiagnosticCreateOrUpdate(ctx, id, payload, diagnostic.DefaultWorkspaceDiagnosticCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DiagnosticClient.WorkspaceDiagnosticDelete`

```go
ctx := context.TODO()
id := diagnostic.NewWorkspaceDiagnosticID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId", "diagnosticId")

read, err := client.WorkspaceDiagnosticDelete(ctx, id, diagnostic.DefaultWorkspaceDiagnosticDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DiagnosticClient.WorkspaceDiagnosticGet`

```go
ctx := context.TODO()
id := diagnostic.NewWorkspaceDiagnosticID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId", "diagnosticId")

read, err := client.WorkspaceDiagnosticGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DiagnosticClient.WorkspaceDiagnosticGetEntityTag`

```go
ctx := context.TODO()
id := diagnostic.NewWorkspaceDiagnosticID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId", "diagnosticId")

read, err := client.WorkspaceDiagnosticGetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DiagnosticClient.WorkspaceDiagnosticListByWorkspace`

```go
ctx := context.TODO()
id := diagnostic.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId")

// alternatively `client.WorkspaceDiagnosticListByWorkspace(ctx, id, diagnostic.DefaultWorkspaceDiagnosticListByWorkspaceOperationOptions())` can be used to do batched pagination
items, err := client.WorkspaceDiagnosticListByWorkspaceComplete(ctx, id, diagnostic.DefaultWorkspaceDiagnosticListByWorkspaceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DiagnosticClient.WorkspaceDiagnosticUpdate`

```go
ctx := context.TODO()
id := diagnostic.NewWorkspaceDiagnosticID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId", "diagnosticId")

payload := diagnostic.DiagnosticUpdateContract{
	// ...
}


read, err := client.WorkspaceDiagnosticUpdate(ctx, id, payload, diagnostic.DefaultWorkspaceDiagnosticUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

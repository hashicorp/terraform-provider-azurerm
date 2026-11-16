
## `github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/apidefinitions` Documentation

The `apidefinitions` SDK allows for interaction with Azure Resource Manager `apicenter` (API Version `2024-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/apidefinitions"
```


### Client Initialization

```go
client := apidefinitions.NewApiDefinitionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ApiDefinitionsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := apidefinitions.NewDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceName", "apiName", "versionName", "definitionName")

payload := apidefinitions.ApiDefinition{
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


### Example Usage: `ApiDefinitionsClient.Delete`

```go
ctx := context.TODO()
id := apidefinitions.NewDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceName", "apiName", "versionName", "definitionName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiDefinitionsClient.ExportSpecification`

```go
ctx := context.TODO()
id := apidefinitions.NewDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceName", "apiName", "versionName", "definitionName")

if err := client.ExportSpecificationThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ApiDefinitionsClient.Get`

```go
ctx := context.TODO()
id := apidefinitions.NewDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceName", "apiName", "versionName", "definitionName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiDefinitionsClient.Head`

```go
ctx := context.TODO()
id := apidefinitions.NewDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceName", "apiName", "versionName", "definitionName")

read, err := client.Head(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiDefinitionsClient.ImportSpecification`

```go
ctx := context.TODO()
id := apidefinitions.NewDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceName", "apiName", "versionName", "definitionName")

payload := apidefinitions.ApiSpecImportRequest{
	// ...
}


if err := client.ImportSpecificationThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ApiDefinitionsClient.List`

```go
ctx := context.TODO()
id := apidefinitions.NewVersionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceName", "apiName", "versionName")

// alternatively `client.List(ctx, id, apidefinitions.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, apidefinitions.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

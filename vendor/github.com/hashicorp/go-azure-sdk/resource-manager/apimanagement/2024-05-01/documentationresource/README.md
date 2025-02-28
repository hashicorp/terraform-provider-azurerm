
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/documentationresource` Documentation

The `documentationresource` SDK allows for interaction with Azure Resource Manager `apimanagement` (API Version `2024-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/documentationresource"
```


### Client Initialization

```go
client := documentationresource.NewDocumentationResourceClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DocumentationResourceClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := documentationresource.NewDocumentationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "documentationId")

payload := documentationresource.DocumentationContract{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, documentationresource.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DocumentationResourceClient.Delete`

```go
ctx := context.TODO()
id := documentationresource.NewDocumentationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "documentationId")

read, err := client.Delete(ctx, id, documentationresource.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DocumentationResourceClient.Get`

```go
ctx := context.TODO()
id := documentationresource.NewDocumentationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "documentationId")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DocumentationResourceClient.GetEntityTag`

```go
ctx := context.TODO()
id := documentationresource.NewDocumentationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "documentationId")

read, err := client.GetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DocumentationResourceClient.ListByService`

```go
ctx := context.TODO()
id := documentationresource.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName")

// alternatively `client.ListByService(ctx, id, documentationresource.DefaultListByServiceOperationOptions())` can be used to do batched pagination
items, err := client.ListByServiceComplete(ctx, id, documentationresource.DefaultListByServiceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DocumentationResourceClient.Update`

```go
ctx := context.TODO()
id := documentationresource.NewDocumentationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "documentationId")

payload := documentationresource.DocumentationUpdateContract{
	// ...
}


read, err := client.Update(ctx, id, payload, documentationresource.DefaultUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

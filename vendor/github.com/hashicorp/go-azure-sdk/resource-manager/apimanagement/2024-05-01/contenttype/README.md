
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/contenttype` Documentation

The `contenttype` SDK allows for interaction with Azure Resource Manager `apimanagement` (API Version `2024-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/contenttype"
```


### Client Initialization

```go
client := contenttype.NewContentTypeClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ContentTypeClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := contenttype.NewContentTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "contentTypeId")

payload := contenttype.ContentTypeContract{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, contenttype.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ContentTypeClient.Delete`

```go
ctx := context.TODO()
id := contenttype.NewContentTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "contentTypeId")

read, err := client.Delete(ctx, id, contenttype.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ContentTypeClient.Get`

```go
ctx := context.TODO()
id := contenttype.NewContentTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "contentTypeId")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ContentTypeClient.ListByService`

```go
ctx := context.TODO()
id := contenttype.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName")

// alternatively `client.ListByService(ctx, id)` can be used to do batched pagination
items, err := client.ListByServiceComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

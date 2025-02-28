
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/productapilink` Documentation

The `productapilink` SDK allows for interaction with Azure Resource Manager `apimanagement` (API Version `2024-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/productapilink"
```


### Client Initialization

```go
client := productapilink.NewProductApiLinkClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ProductApiLinkClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := productapilink.NewProductApiLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "productId", "apiLinkId")

payload := productapilink.ProductApiLinkContract{
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


### Example Usage: `ProductApiLinkClient.Delete`

```go
ctx := context.TODO()
id := productapilink.NewProductApiLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "productId", "apiLinkId")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProductApiLinkClient.Get`

```go
ctx := context.TODO()
id := productapilink.NewProductApiLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "productId", "apiLinkId")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProductApiLinkClient.ListByProduct`

```go
ctx := context.TODO()
id := productapilink.NewProductID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "productId")

// alternatively `client.ListByProduct(ctx, id, productapilink.DefaultListByProductOperationOptions())` can be used to do batched pagination
items, err := client.ListByProductComplete(ctx, id, productapilink.DefaultListByProductOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ProductApiLinkClient.WorkspaceProductApiLinkCreateOrUpdate`

```go
ctx := context.TODO()
id := productapilink.NewWorkspaceProductApiLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId", "productId", "apiLinkId")

payload := productapilink.ProductApiLinkContract{
	// ...
}


read, err := client.WorkspaceProductApiLinkCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProductApiLinkClient.WorkspaceProductApiLinkDelete`

```go
ctx := context.TODO()
id := productapilink.NewWorkspaceProductApiLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId", "productId", "apiLinkId")

read, err := client.WorkspaceProductApiLinkDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProductApiLinkClient.WorkspaceProductApiLinkGet`

```go
ctx := context.TODO()
id := productapilink.NewWorkspaceProductApiLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId", "productId", "apiLinkId")

read, err := client.WorkspaceProductApiLinkGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProductApiLinkClient.WorkspaceProductApiLinkListByProduct`

```go
ctx := context.TODO()
id := productapilink.NewWorkspaceProductID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId", "productId")

// alternatively `client.WorkspaceProductApiLinkListByProduct(ctx, id, productapilink.DefaultWorkspaceProductApiLinkListByProductOperationOptions())` can be used to do batched pagination
items, err := client.WorkspaceProductApiLinkListByProductComplete(ctx, id, productapilink.DefaultWorkspaceProductApiLinkListByProductOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2023-05-01-preview/product` Documentation

The `product` SDK allows for interaction with Azure Resource Manager `apimanagement` (API Version `2023-05-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2023-05-01-preview/product"
```


### Client Initialization

```go
client := product.NewProductClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ProductClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := product.NewProductID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "productId")

payload := product.ProductContract{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, product.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProductClient.Delete`

```go
ctx := context.TODO()
id := product.NewProductID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "productId")

read, err := client.Delete(ctx, id, product.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProductClient.Get`

```go
ctx := context.TODO()
id := product.NewProductID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "productId")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProductClient.GetEntityTag`

```go
ctx := context.TODO()
id := product.NewProductID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "productId")

read, err := client.GetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProductClient.ListByService`

```go
ctx := context.TODO()
id := product.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName")

// alternatively `client.ListByService(ctx, id, product.DefaultListByServiceOperationOptions())` can be used to do batched pagination
items, err := client.ListByServiceComplete(ctx, id, product.DefaultListByServiceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ProductClient.Update`

```go
ctx := context.TODO()
id := product.NewProductID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "productId")

payload := product.ProductUpdateParameters{
	// ...
}


read, err := client.Update(ctx, id, payload, product.DefaultUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProductClient.WorkspaceProductCreateOrUpdate`

```go
ctx := context.TODO()
id := product.NewWorkspaceProductID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId", "productId")

payload := product.ProductContract{
	// ...
}


read, err := client.WorkspaceProductCreateOrUpdate(ctx, id, payload, product.DefaultWorkspaceProductCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProductClient.WorkspaceProductDelete`

```go
ctx := context.TODO()
id := product.NewWorkspaceProductID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId", "productId")

read, err := client.WorkspaceProductDelete(ctx, id, product.DefaultWorkspaceProductDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProductClient.WorkspaceProductGet`

```go
ctx := context.TODO()
id := product.NewWorkspaceProductID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId", "productId")

read, err := client.WorkspaceProductGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProductClient.WorkspaceProductGetEntityTag`

```go
ctx := context.TODO()
id := product.NewWorkspaceProductID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId", "productId")

read, err := client.WorkspaceProductGetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProductClient.WorkspaceProductListByService`

```go
ctx := context.TODO()
id := product.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId")

// alternatively `client.WorkspaceProductListByService(ctx, id, product.DefaultWorkspaceProductListByServiceOperationOptions())` can be used to do batched pagination
items, err := client.WorkspaceProductListByServiceComplete(ctx, id, product.DefaultWorkspaceProductListByServiceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ProductClient.WorkspaceProductUpdate`

```go
ctx := context.TODO()
id := product.NewWorkspaceProductID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId", "productId")

payload := product.ProductUpdateParameters{
	// ...
}


read, err := client.WorkspaceProductUpdate(ctx, id, payload, product.DefaultWorkspaceProductUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

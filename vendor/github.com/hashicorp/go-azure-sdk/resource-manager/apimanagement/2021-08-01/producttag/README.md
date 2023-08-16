
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/producttag` Documentation

The `producttag` SDK allows for interaction with the Azure Resource Manager Service `apimanagement` (API Version `2021-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/producttag"
```


### Client Initialization

```go
client := producttag.NewProductTagClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ProductTagClient.TagAssignToProduct`

```go
ctx := context.TODO()
id := producttag.NewProductTagID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "productIdValue", "tagIdValue")

read, err := client.TagAssignToProduct(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProductTagClient.TagDetachFromProduct`

```go
ctx := context.TODO()
id := producttag.NewProductTagID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "productIdValue", "tagIdValue")

read, err := client.TagDetachFromProduct(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProductTagClient.TagGetByProduct`

```go
ctx := context.TODO()
id := producttag.NewProductTagID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "productIdValue", "tagIdValue")

read, err := client.TagGetByProduct(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProductTagClient.TagGetEntityStateByProduct`

```go
ctx := context.TODO()
id := producttag.NewProductTagID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "productIdValue", "tagIdValue")

read, err := client.TagGetEntityStateByProduct(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProductTagClient.TagListByProduct`

```go
ctx := context.TODO()
id := producttag.NewProductID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "productIdValue")

// alternatively `client.TagListByProduct(ctx, id, producttag.DefaultTagListByProductOperationOptions())` can be used to do batched pagination
items, err := client.TagListByProductComplete(ctx, id, producttag.DefaultTagListByProductOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

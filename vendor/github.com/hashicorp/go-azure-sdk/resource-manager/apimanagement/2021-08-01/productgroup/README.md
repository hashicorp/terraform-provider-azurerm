
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/productgroup` Documentation

The `productgroup` SDK allows for interaction with the Azure Resource Manager Service `apimanagement` (API Version `2021-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/productgroup"
```


### Client Initialization

```go
client := productgroup.NewProductGroupClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ProductGroupClient.CheckEntityExists`

```go
ctx := context.TODO()
id := productgroup.NewProductGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "productIdValue", "groupIdValue")

read, err := client.CheckEntityExists(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProductGroupClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := productgroup.NewProductGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "productIdValue", "groupIdValue")

read, err := client.CreateOrUpdate(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProductGroupClient.Delete`

```go
ctx := context.TODO()
id := productgroup.NewProductGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "productIdValue", "groupIdValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProductGroupClient.ListByProduct`

```go
ctx := context.TODO()
id := productgroup.NewProductID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "productIdValue")

// alternatively `client.ListByProduct(ctx, id, productgroup.DefaultListByProductOperationOptions())` can be used to do batched pagination
items, err := client.ListByProductComplete(ctx, id, productgroup.DefaultListByProductOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

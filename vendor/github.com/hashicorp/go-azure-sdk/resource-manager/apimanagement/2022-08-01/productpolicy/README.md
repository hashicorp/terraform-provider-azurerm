
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/productpolicy` Documentation

The `productpolicy` SDK allows for interaction with Azure Resource Manager `apimanagement` (API Version `2022-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/productpolicy"
```


### Client Initialization

```go
client := productpolicy.NewProductPolicyClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ProductPolicyClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := productpolicy.NewProductID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "productId")

payload := productpolicy.PolicyContract{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, productpolicy.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProductPolicyClient.Delete`

```go
ctx := context.TODO()
id := productpolicy.NewProductID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "productId")

read, err := client.Delete(ctx, id, productpolicy.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProductPolicyClient.Get`

```go
ctx := context.TODO()
id := productpolicy.NewProductID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "productId")

read, err := client.Get(ctx, id, productpolicy.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProductPolicyClient.GetEntityTag`

```go
ctx := context.TODO()
id := productpolicy.NewProductID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "productId")

read, err := client.GetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProductPolicyClient.ListByProduct`

```go
ctx := context.TODO()
id := productpolicy.NewProductID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "productId")

// alternatively `client.ListByProduct(ctx, id)` can be used to do batched pagination
items, err := client.ListByProductComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/productapi` Documentation

The `productapi` SDK allows for interaction with the Azure Resource Manager Service `apimanagement` (API Version `2021-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/productapi"
```


### Client Initialization

```go
client := productapi.NewProductApiClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ProductApiClient.CheckEntityExists`

```go
ctx := context.TODO()
id := productapi.NewProductApiID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "productIdValue", "apiIdValue")

read, err := client.CheckEntityExists(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProductApiClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := productapi.NewProductApiID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "productIdValue", "apiIdValue")

read, err := client.CreateOrUpdate(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProductApiClient.Delete`

```go
ctx := context.TODO()
id := productapi.NewProductApiID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "productIdValue", "apiIdValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProductApiClient.ListByProduct`

```go
ctx := context.TODO()
id := productapi.NewProductID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "productIdValue")

// alternatively `client.ListByProduct(ctx, id, productapi.DefaultListByProductOperationOptions())` can be used to do batched pagination
items, err := client.ListByProductComplete(ctx, id, productapi.DefaultListByProductOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

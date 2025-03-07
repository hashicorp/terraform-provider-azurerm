
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/productgrouplink` Documentation

The `productgrouplink` SDK allows for interaction with Azure Resource Manager `apimanagement` (API Version `2024-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/productgrouplink"
```


### Client Initialization

```go
client := productgrouplink.NewProductGroupLinkClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ProductGroupLinkClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := productgrouplink.NewGroupLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "productId", "groupLinkId")

payload := productgrouplink.ProductGroupLinkContract{
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


### Example Usage: `ProductGroupLinkClient.Delete`

```go
ctx := context.TODO()
id := productgrouplink.NewGroupLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "productId", "groupLinkId")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProductGroupLinkClient.Get`

```go
ctx := context.TODO()
id := productgrouplink.NewGroupLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "productId", "groupLinkId")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProductGroupLinkClient.ListByProduct`

```go
ctx := context.TODO()
id := productgrouplink.NewProductID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "productId")

// alternatively `client.ListByProduct(ctx, id, productgrouplink.DefaultListByProductOperationOptions())` can be used to do batched pagination
items, err := client.ListByProductComplete(ctx, id, productgrouplink.DefaultListByProductOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ProductGroupLinkClient.WorkspaceProductGroupLinkCreateOrUpdate`

```go
ctx := context.TODO()
id := productgrouplink.NewProductGroupLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId", "productId", "groupLinkId")

payload := productgrouplink.ProductGroupLinkContract{
	// ...
}


read, err := client.WorkspaceProductGroupLinkCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProductGroupLinkClient.WorkspaceProductGroupLinkDelete`

```go
ctx := context.TODO()
id := productgrouplink.NewProductGroupLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId", "productId", "groupLinkId")

read, err := client.WorkspaceProductGroupLinkDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProductGroupLinkClient.WorkspaceProductGroupLinkGet`

```go
ctx := context.TODO()
id := productgrouplink.NewProductGroupLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId", "productId", "groupLinkId")

read, err := client.WorkspaceProductGroupLinkGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProductGroupLinkClient.WorkspaceProductGroupLinkListByProduct`

```go
ctx := context.TODO()
id := productgrouplink.NewWorkspaceProductID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId", "productId")

// alternatively `client.WorkspaceProductGroupLinkListByProduct(ctx, id, productgrouplink.DefaultWorkspaceProductGroupLinkListByProductOperationOptions())` can be used to do batched pagination
items, err := client.WorkspaceProductGroupLinkListByProductComplete(ctx, id, productgrouplink.DefaultWorkspaceProductGroupLinkListByProductOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/tagproductlink` Documentation

The `tagproductlink` SDK allows for interaction with Azure Resource Manager `apimanagement` (API Version `2024-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/tagproductlink"
```


### Client Initialization

```go
client := tagproductlink.NewTagProductLinkClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `TagProductLinkClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := tagproductlink.NewProductLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "tagId", "productLinkId")

payload := tagproductlink.TagProductLinkContract{
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


### Example Usage: `TagProductLinkClient.Delete`

```go
ctx := context.TODO()
id := tagproductlink.NewProductLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "tagId", "productLinkId")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TagProductLinkClient.Get`

```go
ctx := context.TODO()
id := tagproductlink.NewProductLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "tagId", "productLinkId")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TagProductLinkClient.ListByProduct`

```go
ctx := context.TODO()
id := tagproductlink.NewTagID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "tagId")

// alternatively `client.ListByProduct(ctx, id, tagproductlink.DefaultListByProductOperationOptions())` can be used to do batched pagination
items, err := client.ListByProductComplete(ctx, id, tagproductlink.DefaultListByProductOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `TagProductLinkClient.WorkspaceTagProductLinkCreateOrUpdate`

```go
ctx := context.TODO()
id := tagproductlink.NewTagProductLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId", "tagId", "productLinkId")

payload := tagproductlink.TagProductLinkContract{
	// ...
}


read, err := client.WorkspaceTagProductLinkCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TagProductLinkClient.WorkspaceTagProductLinkDelete`

```go
ctx := context.TODO()
id := tagproductlink.NewTagProductLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId", "tagId", "productLinkId")

read, err := client.WorkspaceTagProductLinkDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TagProductLinkClient.WorkspaceTagProductLinkGet`

```go
ctx := context.TODO()
id := tagproductlink.NewTagProductLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId", "tagId", "productLinkId")

read, err := client.WorkspaceTagProductLinkGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TagProductLinkClient.WorkspaceTagProductLinkListByProduct`

```go
ctx := context.TODO()
id := tagproductlink.NewWorkspaceTagID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId", "tagId")

// alternatively `client.WorkspaceTagProductLinkListByProduct(ctx, id, tagproductlink.DefaultWorkspaceTagProductLinkListByProductOperationOptions())` can be used to do batched pagination
items, err := client.WorkspaceTagProductLinkListByProductComplete(ctx, id, tagproductlink.DefaultWorkspaceTagProductLinkListByProductOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

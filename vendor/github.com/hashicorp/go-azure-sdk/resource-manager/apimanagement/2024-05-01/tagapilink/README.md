
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/tagapilink` Documentation

The `tagapilink` SDK allows for interaction with Azure Resource Manager `apimanagement` (API Version `2024-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/tagapilink"
```


### Client Initialization

```go
client := tagapilink.NewTagApiLinkClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `TagApiLinkClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := tagapilink.NewApiLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "tagId", "apiLinkId")

payload := tagapilink.TagApiLinkContract{
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


### Example Usage: `TagApiLinkClient.Delete`

```go
ctx := context.TODO()
id := tagapilink.NewApiLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "tagId", "apiLinkId")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TagApiLinkClient.Get`

```go
ctx := context.TODO()
id := tagapilink.NewApiLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "tagId", "apiLinkId")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TagApiLinkClient.ListByProduct`

```go
ctx := context.TODO()
id := tagapilink.NewTagID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "tagId")

// alternatively `client.ListByProduct(ctx, id, tagapilink.DefaultListByProductOperationOptions())` can be used to do batched pagination
items, err := client.ListByProductComplete(ctx, id, tagapilink.DefaultListByProductOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `TagApiLinkClient.WorkspaceTagApiLinkCreateOrUpdate`

```go
ctx := context.TODO()
id := tagapilink.NewTagApiLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId", "tagId", "apiLinkId")

payload := tagapilink.TagApiLinkContract{
	// ...
}


read, err := client.WorkspaceTagApiLinkCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TagApiLinkClient.WorkspaceTagApiLinkDelete`

```go
ctx := context.TODO()
id := tagapilink.NewTagApiLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId", "tagId", "apiLinkId")

read, err := client.WorkspaceTagApiLinkDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TagApiLinkClient.WorkspaceTagApiLinkGet`

```go
ctx := context.TODO()
id := tagapilink.NewTagApiLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId", "tagId", "apiLinkId")

read, err := client.WorkspaceTagApiLinkGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TagApiLinkClient.WorkspaceTagApiLinkListByProduct`

```go
ctx := context.TODO()
id := tagapilink.NewWorkspaceTagID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId", "tagId")

// alternatively `client.WorkspaceTagApiLinkListByProduct(ctx, id, tagapilink.DefaultWorkspaceTagApiLinkListByProductOperationOptions())` can be used to do batched pagination
items, err := client.WorkspaceTagApiLinkListByProductComplete(ctx, id, tagapilink.DefaultWorkspaceTagApiLinkListByProductOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

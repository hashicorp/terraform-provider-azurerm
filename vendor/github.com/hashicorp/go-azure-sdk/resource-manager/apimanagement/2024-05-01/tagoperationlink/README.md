
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/tagoperationlink` Documentation

The `tagoperationlink` SDK allows for interaction with Azure Resource Manager `apimanagement` (API Version `2024-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/tagoperationlink"
```


### Client Initialization

```go
client := tagoperationlink.NewTagOperationLinkClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `TagOperationLinkClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := tagoperationlink.NewOperationLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "tagId", "operationLinkId")

payload := tagoperationlink.TagOperationLinkContract{
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


### Example Usage: `TagOperationLinkClient.Delete`

```go
ctx := context.TODO()
id := tagoperationlink.NewOperationLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "tagId", "operationLinkId")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TagOperationLinkClient.Get`

```go
ctx := context.TODO()
id := tagoperationlink.NewOperationLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "tagId", "operationLinkId")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TagOperationLinkClient.ListByProduct`

```go
ctx := context.TODO()
id := tagoperationlink.NewTagID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "tagId")

// alternatively `client.ListByProduct(ctx, id, tagoperationlink.DefaultListByProductOperationOptions())` can be used to do batched pagination
items, err := client.ListByProductComplete(ctx, id, tagoperationlink.DefaultListByProductOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `TagOperationLinkClient.WorkspaceTagOperationLinkCreateOrUpdate`

```go
ctx := context.TODO()
id := tagoperationlink.NewTagOperationLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId", "tagId", "operationLinkId")

payload := tagoperationlink.TagOperationLinkContract{
	// ...
}


read, err := client.WorkspaceTagOperationLinkCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TagOperationLinkClient.WorkspaceTagOperationLinkDelete`

```go
ctx := context.TODO()
id := tagoperationlink.NewTagOperationLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId", "tagId", "operationLinkId")

read, err := client.WorkspaceTagOperationLinkDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TagOperationLinkClient.WorkspaceTagOperationLinkGet`

```go
ctx := context.TODO()
id := tagoperationlink.NewTagOperationLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId", "tagId", "operationLinkId")

read, err := client.WorkspaceTagOperationLinkGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TagOperationLinkClient.WorkspaceTagOperationLinkListByProduct`

```go
ctx := context.TODO()
id := tagoperationlink.NewWorkspaceTagID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId", "tagId")

// alternatively `client.WorkspaceTagOperationLinkListByProduct(ctx, id, tagoperationlink.DefaultWorkspaceTagOperationLinkListByProductOperationOptions())` can be used to do batched pagination
items, err := client.WorkspaceTagOperationLinkListByProductComplete(ctx, id, tagoperationlink.DefaultWorkspaceTagOperationLinkListByProductOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

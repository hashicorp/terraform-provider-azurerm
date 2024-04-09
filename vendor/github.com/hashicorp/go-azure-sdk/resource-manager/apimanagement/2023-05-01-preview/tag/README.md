
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2023-05-01-preview/tag` Documentation

The `tag` SDK allows for interaction with the Azure Resource Manager Service `apimanagement` (API Version `2023-05-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2023-05-01-preview/tag"
```


### Client Initialization

```go
client := tag.NewTagClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `TagClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := tag.NewTagID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "tagIdValue")

payload := tag.TagCreateUpdateParameters{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, tag.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TagClient.Delete`

```go
ctx := context.TODO()
id := tag.NewTagID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "tagIdValue")

read, err := client.Delete(ctx, id, tag.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TagClient.Get`

```go
ctx := context.TODO()
id := tag.NewTagID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "tagIdValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TagClient.GetEntityState`

```go
ctx := context.TODO()
id := tag.NewTagID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "tagIdValue")

read, err := client.GetEntityState(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TagClient.ListByService`

```go
ctx := context.TODO()
id := tag.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue")

// alternatively `client.ListByService(ctx, id, tag.DefaultListByServiceOperationOptions())` can be used to do batched pagination
items, err := client.ListByServiceComplete(ctx, id, tag.DefaultListByServiceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `TagClient.Update`

```go
ctx := context.TODO()
id := tag.NewTagID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "tagIdValue")

payload := tag.TagCreateUpdateParameters{
	// ...
}


read, err := client.Update(ctx, id, payload, tag.DefaultUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TagClient.WorkspaceTagCreateOrUpdate`

```go
ctx := context.TODO()
id := tag.NewWorkspaceTagID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "tagIdValue")

payload := tag.TagCreateUpdateParameters{
	// ...
}


read, err := client.WorkspaceTagCreateOrUpdate(ctx, id, payload, tag.DefaultWorkspaceTagCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TagClient.WorkspaceTagDelete`

```go
ctx := context.TODO()
id := tag.NewWorkspaceTagID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "tagIdValue")

read, err := client.WorkspaceTagDelete(ctx, id, tag.DefaultWorkspaceTagDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TagClient.WorkspaceTagGet`

```go
ctx := context.TODO()
id := tag.NewWorkspaceTagID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "tagIdValue")

read, err := client.WorkspaceTagGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TagClient.WorkspaceTagGetEntityState`

```go
ctx := context.TODO()
id := tag.NewWorkspaceTagID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "tagIdValue")

read, err := client.WorkspaceTagGetEntityState(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TagClient.WorkspaceTagListByService`

```go
ctx := context.TODO()
id := tag.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue")

// alternatively `client.WorkspaceTagListByService(ctx, id, tag.DefaultWorkspaceTagListByServiceOperationOptions())` can be used to do batched pagination
items, err := client.WorkspaceTagListByServiceComplete(ctx, id, tag.DefaultWorkspaceTagListByServiceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `TagClient.WorkspaceTagUpdate`

```go
ctx := context.TODO()
id := tag.NewWorkspaceTagID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "tagIdValue")

payload := tag.TagCreateUpdateParameters{
	// ...
}


read, err := client.WorkspaceTagUpdate(ctx, id, payload, tag.DefaultWorkspaceTagUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

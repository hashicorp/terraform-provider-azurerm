
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2023-05-01-preview/group` Documentation

The `group` SDK allows for interaction with the Azure Resource Manager Service `apimanagement` (API Version `2023-05-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2023-05-01-preview/group"
```


### Client Initialization

```go
client := group.NewGroupClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `GroupClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := group.NewGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "groupIdValue")

payload := group.GroupCreateParameters{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, group.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GroupClient.Delete`

```go
ctx := context.TODO()
id := group.NewGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "groupIdValue")

read, err := client.Delete(ctx, id, group.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GroupClient.Get`

```go
ctx := context.TODO()
id := group.NewGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "groupIdValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GroupClient.GetEntityTag`

```go
ctx := context.TODO()
id := group.NewGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "groupIdValue")

read, err := client.GetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GroupClient.ListByService`

```go
ctx := context.TODO()
id := group.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue")

// alternatively `client.ListByService(ctx, id, group.DefaultListByServiceOperationOptions())` can be used to do batched pagination
items, err := client.ListByServiceComplete(ctx, id, group.DefaultListByServiceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `GroupClient.Update`

```go
ctx := context.TODO()
id := group.NewGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "groupIdValue")

payload := group.GroupUpdateParameters{
	// ...
}


read, err := client.Update(ctx, id, payload, group.DefaultUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GroupClient.WorkspaceGroupCreateOrUpdate`

```go
ctx := context.TODO()
id := group.NewWorkspaceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "groupIdValue")

payload := group.GroupCreateParameters{
	// ...
}


read, err := client.WorkspaceGroupCreateOrUpdate(ctx, id, payload, group.DefaultWorkspaceGroupCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GroupClient.WorkspaceGroupDelete`

```go
ctx := context.TODO()
id := group.NewWorkspaceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "groupIdValue")

read, err := client.WorkspaceGroupDelete(ctx, id, group.DefaultWorkspaceGroupDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GroupClient.WorkspaceGroupGet`

```go
ctx := context.TODO()
id := group.NewWorkspaceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "groupIdValue")

read, err := client.WorkspaceGroupGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GroupClient.WorkspaceGroupGetEntityTag`

```go
ctx := context.TODO()
id := group.NewWorkspaceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "groupIdValue")

read, err := client.WorkspaceGroupGetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GroupClient.WorkspaceGroupListByService`

```go
ctx := context.TODO()
id := group.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue")

// alternatively `client.WorkspaceGroupListByService(ctx, id, group.DefaultWorkspaceGroupListByServiceOperationOptions())` can be used to do batched pagination
items, err := client.WorkspaceGroupListByServiceComplete(ctx, id, group.DefaultWorkspaceGroupListByServiceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `GroupClient.WorkspaceGroupUpdate`

```go
ctx := context.TODO()
id := group.NewWorkspaceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "groupIdValue")

payload := group.GroupUpdateParameters{
	// ...
}


read, err := client.WorkspaceGroupUpdate(ctx, id, payload, group.DefaultWorkspaceGroupUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2023-05-01-preview/api` Documentation

The `api` SDK allows for interaction with the Azure Resource Manager Service `apimanagement` (API Version `2023-05-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2023-05-01-preview/api"
```


### Client Initialization

```go
client := api.NewApiClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ApiClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := api.NewApiID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "apiIdValue")

payload := api.ApiCreateOrUpdateParameter{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload, api.DefaultCreateOrUpdateOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `ApiClient.Delete`

```go
ctx := context.TODO()
id := api.NewApiID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "apiIdValue")

if err := client.DeleteThenPoll(ctx, id, api.DefaultDeleteOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `ApiClient.Get`

```go
ctx := context.TODO()
id := api.NewApiID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "apiIdValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiClient.GetEntityTag`

```go
ctx := context.TODO()
id := api.NewApiID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "apiIdValue")

read, err := client.GetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiClient.ListByService`

```go
ctx := context.TODO()
id := api.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue")

// alternatively `client.ListByService(ctx, id, api.DefaultListByServiceOperationOptions())` can be used to do batched pagination
items, err := client.ListByServiceComplete(ctx, id, api.DefaultListByServiceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ApiClient.Update`

```go
ctx := context.TODO()
id := api.NewApiID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "apiIdValue")

payload := api.ApiUpdateContract{
	// ...
}


read, err := client.Update(ctx, id, payload, api.DefaultUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiClient.WorkspaceApiCreateOrUpdate`

```go
ctx := context.TODO()
id := api.NewWorkspaceApiID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "apiIdValue")

payload := api.ApiCreateOrUpdateParameter{
	// ...
}


if err := client.WorkspaceApiCreateOrUpdateThenPoll(ctx, id, payload, api.DefaultWorkspaceApiCreateOrUpdateOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `ApiClient.WorkspaceApiDelete`

```go
ctx := context.TODO()
id := api.NewWorkspaceApiID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "apiIdValue")

read, err := client.WorkspaceApiDelete(ctx, id, api.DefaultWorkspaceApiDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiClient.WorkspaceApiGet`

```go
ctx := context.TODO()
id := api.NewWorkspaceApiID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "apiIdValue")

read, err := client.WorkspaceApiGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiClient.WorkspaceApiGetEntityTag`

```go
ctx := context.TODO()
id := api.NewWorkspaceApiID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "apiIdValue")

read, err := client.WorkspaceApiGetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiClient.WorkspaceApiListByService`

```go
ctx := context.TODO()
id := api.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue")

// alternatively `client.WorkspaceApiListByService(ctx, id, api.DefaultWorkspaceApiListByServiceOperationOptions())` can be used to do batched pagination
items, err := client.WorkspaceApiListByServiceComplete(ctx, id, api.DefaultWorkspaceApiListByServiceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ApiClient.WorkspaceApiUpdate`

```go
ctx := context.TODO()
id := api.NewWorkspaceApiID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "apiIdValue")

payload := api.ApiUpdateContract{
	// ...
}


read, err := client.WorkspaceApiUpdate(ctx, id, payload, api.DefaultWorkspaceApiUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

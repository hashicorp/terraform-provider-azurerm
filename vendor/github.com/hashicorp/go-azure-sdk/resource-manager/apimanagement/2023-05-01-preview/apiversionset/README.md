
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2023-05-01-preview/apiversionset` Documentation

The `apiversionset` SDK allows for interaction with Azure Resource Manager `apimanagement` (API Version `2023-05-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2023-05-01-preview/apiversionset"
```


### Client Initialization

```go
client := apiversionset.NewApiVersionSetClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ApiVersionSetClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := apiversionset.NewApiVersionSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "versionSetId")

payload := apiversionset.ApiVersionSetContract{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, apiversionset.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiVersionSetClient.Get`

```go
ctx := context.TODO()
id := apiversionset.NewApiVersionSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "versionSetId")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiVersionSetClient.GetEntityTag`

```go
ctx := context.TODO()
id := apiversionset.NewApiVersionSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "versionSetId")

read, err := client.GetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiVersionSetClient.ListByService`

```go
ctx := context.TODO()
id := apiversionset.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName")

// alternatively `client.ListByService(ctx, id, apiversionset.DefaultListByServiceOperationOptions())` can be used to do batched pagination
items, err := client.ListByServiceComplete(ctx, id, apiversionset.DefaultListByServiceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ApiVersionSetClient.Update`

```go
ctx := context.TODO()
id := apiversionset.NewApiVersionSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "versionSetId")

payload := apiversionset.ApiVersionSetUpdateParameters{
	// ...
}


read, err := client.Update(ctx, id, payload, apiversionset.DefaultUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiVersionSetClient.WorkspaceApiVersionSetCreateOrUpdate`

```go
ctx := context.TODO()
id := apiversionset.NewWorkspaceApiVersionSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId", "versionSetId")

payload := apiversionset.ApiVersionSetContract{
	// ...
}


read, err := client.WorkspaceApiVersionSetCreateOrUpdate(ctx, id, payload, apiversionset.DefaultWorkspaceApiVersionSetCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiVersionSetClient.WorkspaceApiVersionSetGet`

```go
ctx := context.TODO()
id := apiversionset.NewWorkspaceApiVersionSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId", "versionSetId")

read, err := client.WorkspaceApiVersionSetGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiVersionSetClient.WorkspaceApiVersionSetGetEntityTag`

```go
ctx := context.TODO()
id := apiversionset.NewWorkspaceApiVersionSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId", "versionSetId")

read, err := client.WorkspaceApiVersionSetGetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiVersionSetClient.WorkspaceApiVersionSetListByService`

```go
ctx := context.TODO()
id := apiversionset.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId")

// alternatively `client.WorkspaceApiVersionSetListByService(ctx, id, apiversionset.DefaultWorkspaceApiVersionSetListByServiceOperationOptions())` can be used to do batched pagination
items, err := client.WorkspaceApiVersionSetListByServiceComplete(ctx, id, apiversionset.DefaultWorkspaceApiVersionSetListByServiceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ApiVersionSetClient.WorkspaceApiVersionSetUpdate`

```go
ctx := context.TODO()
id := apiversionset.NewWorkspaceApiVersionSetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId", "versionSetId")

payload := apiversionset.ApiVersionSetUpdateParameters{
	// ...
}


read, err := client.WorkspaceApiVersionSetUpdate(ctx, id, payload, apiversionset.DefaultWorkspaceApiVersionSetUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

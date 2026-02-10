
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/logger` Documentation

The `logger` SDK allows for interaction with Azure Resource Manager `apimanagement` (API Version `2024-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/logger"
```


### Client Initialization

```go
client := logger.NewLoggerClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `LoggerClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := logger.NewLoggerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "loggerId")

payload := logger.LoggerContract{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, logger.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LoggerClient.Delete`

```go
ctx := context.TODO()
id := logger.NewLoggerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "loggerId")

read, err := client.Delete(ctx, id, logger.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LoggerClient.Get`

```go
ctx := context.TODO()
id := logger.NewLoggerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "loggerId")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LoggerClient.GetEntityTag`

```go
ctx := context.TODO()
id := logger.NewLoggerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "loggerId")

read, err := client.GetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LoggerClient.ListByService`

```go
ctx := context.TODO()
id := logger.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName")

// alternatively `client.ListByService(ctx, id, logger.DefaultListByServiceOperationOptions())` can be used to do batched pagination
items, err := client.ListByServiceComplete(ctx, id, logger.DefaultListByServiceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `LoggerClient.Update`

```go
ctx := context.TODO()
id := logger.NewLoggerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "loggerId")

payload := logger.LoggerUpdateContract{
	// ...
}


read, err := client.Update(ctx, id, payload, logger.DefaultUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LoggerClient.WorkspaceLoggerCreateOrUpdate`

```go
ctx := context.TODO()
id := logger.NewWorkspaceLoggerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId", "loggerId")

payload := logger.LoggerContract{
	// ...
}


read, err := client.WorkspaceLoggerCreateOrUpdate(ctx, id, payload, logger.DefaultWorkspaceLoggerCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LoggerClient.WorkspaceLoggerDelete`

```go
ctx := context.TODO()
id := logger.NewWorkspaceLoggerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId", "loggerId")

read, err := client.WorkspaceLoggerDelete(ctx, id, logger.DefaultWorkspaceLoggerDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LoggerClient.WorkspaceLoggerGet`

```go
ctx := context.TODO()
id := logger.NewWorkspaceLoggerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId", "loggerId")

read, err := client.WorkspaceLoggerGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LoggerClient.WorkspaceLoggerGetEntityTag`

```go
ctx := context.TODO()
id := logger.NewWorkspaceLoggerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId", "loggerId")

read, err := client.WorkspaceLoggerGetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LoggerClient.WorkspaceLoggerListByWorkspace`

```go
ctx := context.TODO()
id := logger.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId")

// alternatively `client.WorkspaceLoggerListByWorkspace(ctx, id, logger.DefaultWorkspaceLoggerListByWorkspaceOperationOptions())` can be used to do batched pagination
items, err := client.WorkspaceLoggerListByWorkspaceComplete(ctx, id, logger.DefaultWorkspaceLoggerListByWorkspaceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `LoggerClient.WorkspaceLoggerUpdate`

```go
ctx := context.TODO()
id := logger.NewWorkspaceLoggerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId", "loggerId")

payload := logger.LoggerUpdateContract{
	// ...
}


read, err := client.WorkspaceLoggerUpdate(ctx, id, payload, logger.DefaultWorkspaceLoggerUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

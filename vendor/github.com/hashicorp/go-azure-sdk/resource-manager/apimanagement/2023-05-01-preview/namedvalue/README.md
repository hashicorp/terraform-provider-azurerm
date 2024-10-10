
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2023-05-01-preview/namedvalue` Documentation

The `namedvalue` SDK allows for interaction with the Azure Resource Manager Service `apimanagement` (API Version `2023-05-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2023-05-01-preview/namedvalue"
```


### Client Initialization

```go
client := namedvalue.NewNamedValueClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `NamedValueClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := namedvalue.NewNamedValueID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "namedValueIdValue")

payload := namedvalue.NamedValueCreateContract{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload, namedvalue.DefaultCreateOrUpdateOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `NamedValueClient.Delete`

```go
ctx := context.TODO()
id := namedvalue.NewNamedValueID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "namedValueIdValue")

read, err := client.Delete(ctx, id, namedvalue.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NamedValueClient.Get`

```go
ctx := context.TODO()
id := namedvalue.NewNamedValueID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "namedValueIdValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NamedValueClient.GetEntityTag`

```go
ctx := context.TODO()
id := namedvalue.NewNamedValueID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "namedValueIdValue")

read, err := client.GetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NamedValueClient.ListByService`

```go
ctx := context.TODO()
id := namedvalue.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue")

// alternatively `client.ListByService(ctx, id, namedvalue.DefaultListByServiceOperationOptions())` can be used to do batched pagination
items, err := client.ListByServiceComplete(ctx, id, namedvalue.DefaultListByServiceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NamedValueClient.ListValue`

```go
ctx := context.TODO()
id := namedvalue.NewNamedValueID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "namedValueIdValue")

read, err := client.ListValue(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NamedValueClient.RefreshSecret`

```go
ctx := context.TODO()
id := namedvalue.NewNamedValueID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "namedValueIdValue")

if err := client.RefreshSecretThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `NamedValueClient.Update`

```go
ctx := context.TODO()
id := namedvalue.NewNamedValueID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "namedValueIdValue")

payload := namedvalue.NamedValueUpdateParameters{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload, namedvalue.DefaultUpdateOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `NamedValueClient.WorkspaceNamedValueCreateOrUpdate`

```go
ctx := context.TODO()
id := namedvalue.NewWorkspaceNamedValueID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "namedValueIdValue")

payload := namedvalue.NamedValueCreateContract{
	// ...
}


if err := client.WorkspaceNamedValueCreateOrUpdateThenPoll(ctx, id, payload, namedvalue.DefaultWorkspaceNamedValueCreateOrUpdateOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `NamedValueClient.WorkspaceNamedValueDelete`

```go
ctx := context.TODO()
id := namedvalue.NewWorkspaceNamedValueID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "namedValueIdValue")

read, err := client.WorkspaceNamedValueDelete(ctx, id, namedvalue.DefaultWorkspaceNamedValueDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NamedValueClient.WorkspaceNamedValueGet`

```go
ctx := context.TODO()
id := namedvalue.NewWorkspaceNamedValueID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "namedValueIdValue")

read, err := client.WorkspaceNamedValueGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NamedValueClient.WorkspaceNamedValueGetEntityTag`

```go
ctx := context.TODO()
id := namedvalue.NewWorkspaceNamedValueID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "namedValueIdValue")

read, err := client.WorkspaceNamedValueGetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NamedValueClient.WorkspaceNamedValueListByService`

```go
ctx := context.TODO()
id := namedvalue.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue")

// alternatively `client.WorkspaceNamedValueListByService(ctx, id, namedvalue.DefaultWorkspaceNamedValueListByServiceOperationOptions())` can be used to do batched pagination
items, err := client.WorkspaceNamedValueListByServiceComplete(ctx, id, namedvalue.DefaultWorkspaceNamedValueListByServiceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NamedValueClient.WorkspaceNamedValueListValue`

```go
ctx := context.TODO()
id := namedvalue.NewWorkspaceNamedValueID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "namedValueIdValue")

read, err := client.WorkspaceNamedValueListValue(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NamedValueClient.WorkspaceNamedValueRefreshSecret`

```go
ctx := context.TODO()
id := namedvalue.NewWorkspaceNamedValueID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "namedValueIdValue")

if err := client.WorkspaceNamedValueRefreshSecretThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `NamedValueClient.WorkspaceNamedValueUpdate`

```go
ctx := context.TODO()
id := namedvalue.NewWorkspaceNamedValueID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "workspaceIdValue", "namedValueIdValue")

payload := namedvalue.NamedValueUpdateParameters{
	// ...
}


if err := client.WorkspaceNamedValueUpdateThenPoll(ctx, id, payload, namedvalue.DefaultWorkspaceNamedValueUpdateOperationOptions()); err != nil {
	// handle the error
}
```

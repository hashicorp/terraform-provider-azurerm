
## `github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/apis` Documentation

The `apis` SDK allows for interaction with Azure Resource Manager `apicenter` (API Version `2024-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/apis"
```


### Client Initialization

```go
client := apis.NewApisClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ApisClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := apis.NewApiID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceName", "apiName")

payload := apis.Api{
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


### Example Usage: `ApisClient.Delete`

```go
ctx := context.TODO()
id := apis.NewApiID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceName", "apiName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApisClient.Get`

```go
ctx := context.TODO()
id := apis.NewApiID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceName", "apiName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApisClient.Head`

```go
ctx := context.TODO()
id := apis.NewApiID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceName", "apiName")

read, err := client.Head(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApisClient.List`

```go
ctx := context.TODO()
id := apis.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceName")

// alternatively `client.List(ctx, id, apis.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, apis.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

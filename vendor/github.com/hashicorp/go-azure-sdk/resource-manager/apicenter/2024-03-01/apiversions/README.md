
## `github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/apiversions` Documentation

The `apiversions` SDK allows for interaction with Azure Resource Manager `apicenter` (API Version `2024-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/apiversions"
```


### Client Initialization

```go
client := apiversions.NewApiVersionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ApiVersionsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := apiversions.NewVersionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceName", "apiName", "versionName")

payload := apiversions.ApiVersion{
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


### Example Usage: `ApiVersionsClient.Delete`

```go
ctx := context.TODO()
id := apiversions.NewVersionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceName", "apiName", "versionName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiVersionsClient.Get`

```go
ctx := context.TODO()
id := apiversions.NewVersionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceName", "apiName", "versionName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiVersionsClient.Head`

```go
ctx := context.TODO()
id := apiversions.NewVersionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceName", "apiName", "versionName")

read, err := client.Head(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiVersionsClient.List`

```go
ctx := context.TODO()
id := apiversions.NewApiID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceName", "apiName")

// alternatively `client.List(ctx, id, apiversions.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, apiversions.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

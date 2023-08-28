
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/apitag` Documentation

The `apitag` SDK allows for interaction with the Azure Resource Manager Service `apimanagement` (API Version `2021-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/apitag"
```


### Client Initialization

```go
client := apitag.NewApiTagClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ApiTagClient.TagAssignToApi`

```go
ctx := context.TODO()
id := apitag.NewApiTagID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "apiIdValue", "tagIdValue")

read, err := client.TagAssignToApi(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiTagClient.TagDetachFromApi`

```go
ctx := context.TODO()
id := apitag.NewApiTagID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "apiIdValue", "tagIdValue")

read, err := client.TagDetachFromApi(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiTagClient.TagGetByApi`

```go
ctx := context.TODO()
id := apitag.NewApiTagID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "apiIdValue", "tagIdValue")

read, err := client.TagGetByApi(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiTagClient.TagGetEntityStateByApi`

```go
ctx := context.TODO()
id := apitag.NewApiTagID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "apiIdValue", "tagIdValue")

read, err := client.TagGetEntityStateByApi(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiTagClient.TagListByApi`

```go
ctx := context.TODO()
id := apitag.NewApiID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "apiIdValue")

// alternatively `client.TagListByApi(ctx, id, apitag.DefaultTagListByApiOperationOptions())` can be used to do batched pagination
items, err := client.TagListByApiComplete(ctx, id, apitag.DefaultTagListByApiOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

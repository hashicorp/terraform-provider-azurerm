
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/apioperationtag` Documentation

The `apioperationtag` SDK allows for interaction with the Azure Resource Manager Service `apimanagement` (API Version `2021-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/apioperationtag"
```


### Client Initialization

```go
client := apioperationtag.NewApiOperationTagClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ApiOperationTagClient.TagAssignToOperation`

```go
ctx := context.TODO()
id := apioperationtag.NewOperationTagID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "apiIdValue", "operationIdValue", "tagIdValue")

read, err := client.TagAssignToOperation(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiOperationTagClient.TagDetachFromOperation`

```go
ctx := context.TODO()
id := apioperationtag.NewOperationTagID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "apiIdValue", "operationIdValue", "tagIdValue")

read, err := client.TagDetachFromOperation(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiOperationTagClient.TagGetByOperation`

```go
ctx := context.TODO()
id := apioperationtag.NewOperationTagID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "apiIdValue", "operationIdValue", "tagIdValue")

read, err := client.TagGetByOperation(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiOperationTagClient.TagGetEntityStateByOperation`

```go
ctx := context.TODO()
id := apioperationtag.NewOperationTagID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "apiIdValue", "operationIdValue", "tagIdValue")

read, err := client.TagGetEntityStateByOperation(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiOperationTagClient.TagListByOperation`

```go
ctx := context.TODO()
id := apioperationtag.NewOperationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "apiIdValue", "operationIdValue")

// alternatively `client.TagListByOperation(ctx, id, apioperationtag.DefaultTagListByOperationOperationOptions())` can be used to do batched pagination
items, err := client.TagListByOperationComplete(ctx, id, apioperationtag.DefaultTagListByOperationOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

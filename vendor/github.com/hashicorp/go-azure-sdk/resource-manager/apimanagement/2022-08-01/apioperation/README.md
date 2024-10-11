
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apioperation` Documentation

The `apioperation` SDK allows for interaction with Azure Resource Manager `apimanagement` (API Version `2022-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apioperation"
```


### Client Initialization

```go
client := apioperation.NewApiOperationClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ApiOperationClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := apioperation.NewOperationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId", "operationId")

payload := apioperation.OperationContract{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, apioperation.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiOperationClient.Delete`

```go
ctx := context.TODO()
id := apioperation.NewOperationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId", "operationId")

read, err := client.Delete(ctx, id, apioperation.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiOperationClient.Get`

```go
ctx := context.TODO()
id := apioperation.NewOperationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId", "operationId")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiOperationClient.GetEntityTag`

```go
ctx := context.TODO()
id := apioperation.NewOperationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId", "operationId")

read, err := client.GetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiOperationClient.ListByApi`

```go
ctx := context.TODO()
id := apioperation.NewApiID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId")

// alternatively `client.ListByApi(ctx, id, apioperation.DefaultListByApiOperationOptions())` can be used to do batched pagination
items, err := client.ListByApiComplete(ctx, id, apioperation.DefaultListByApiOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ApiOperationClient.Update`

```go
ctx := context.TODO()
id := apioperation.NewOperationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId", "operationId")

payload := apioperation.OperationUpdateContract{
	// ...
}


read, err := client.Update(ctx, id, payload, apioperation.DefaultUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

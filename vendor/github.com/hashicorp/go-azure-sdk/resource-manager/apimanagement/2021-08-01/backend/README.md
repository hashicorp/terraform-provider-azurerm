
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/backend` Documentation

The `backend` SDK allows for interaction with the Azure Resource Manager Service `apimanagement` (API Version `2021-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/backend"
```


### Client Initialization

```go
client := backend.NewBackendClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `BackendClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := backend.NewBackendID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "backendIdValue")

payload := backend.BackendContract{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, backend.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BackendClient.Delete`

```go
ctx := context.TODO()
id := backend.NewBackendID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "backendIdValue")

read, err := client.Delete(ctx, id, backend.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BackendClient.Get`

```go
ctx := context.TODO()
id := backend.NewBackendID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "backendIdValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BackendClient.GetEntityTag`

```go
ctx := context.TODO()
id := backend.NewBackendID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "backendIdValue")

read, err := client.GetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BackendClient.ListByService`

```go
ctx := context.TODO()
id := backend.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue")

// alternatively `client.ListByService(ctx, id, backend.DefaultListByServiceOperationOptions())` can be used to do batched pagination
items, err := client.ListByServiceComplete(ctx, id, backend.DefaultListByServiceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `BackendClient.Update`

```go
ctx := context.TODO()
id := backend.NewBackendID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "backendIdValue")

payload := backend.BackendUpdateParameters{
	// ...
}


read, err := client.Update(ctx, id, payload, backend.DefaultUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

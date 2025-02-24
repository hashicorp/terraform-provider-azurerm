
## `github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2023-03-01/endpoints` Documentation

The `endpoints` SDK allows for interaction with Azure Resource Manager `storagemover` (API Version `2023-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2023-03-01/endpoints"
```


### Client Initialization

```go
client := endpoints.NewEndpointsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `EndpointsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := endpoints.NewEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageMoverName", "endpointName")

payload := endpoints.Endpoint{
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


### Example Usage: `EndpointsClient.Delete`

```go
ctx := context.TODO()
id := endpoints.NewEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageMoverName", "endpointName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `EndpointsClient.Get`

```go
ctx := context.TODO()
id := endpoints.NewEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageMoverName", "endpointName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EndpointsClient.List`

```go
ctx := context.TODO()
id := endpoints.NewStorageMoverID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageMoverName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `EndpointsClient.Update`

```go
ctx := context.TODO()
id := endpoints.NewEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageMoverName", "endpointName")

payload := endpoints.EndpointBaseUpdateParameters{
	// ...
}


read, err := client.Update(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

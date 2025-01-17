
## `github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/origins` Documentation

The `origins` SDK allows for interaction with Azure Resource Manager `cdn` (API Version `2024-02-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/origins"
```


### Client Initialization

```go
client := origins.NewOriginsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `OriginsClient.Create`

```go
ctx := context.TODO()
id := origins.NewOriginID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "endpointName", "originName")

payload := origins.Origin{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OriginsClient.Delete`

```go
ctx := context.TODO()
id := origins.NewOriginID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "endpointName", "originName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OriginsClient.Get`

```go
ctx := context.TODO()
id := origins.NewOriginID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "endpointName", "originName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OriginsClient.ListByEndpoint`

```go
ctx := context.TODO()
id := origins.NewEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "endpointName")

// alternatively `client.ListByEndpoint(ctx, id)` can be used to do batched pagination
items, err := client.ListByEndpointComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OriginsClient.Update`

```go
ctx := context.TODO()
id := origins.NewOriginID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileName", "endpointName", "originName")

payload := origins.OriginUpdateParameters{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

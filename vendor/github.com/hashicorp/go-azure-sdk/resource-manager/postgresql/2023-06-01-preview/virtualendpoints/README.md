
## `github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2023-06-01-preview/virtualendpoints` Documentation

The `virtualendpoints` SDK allows for interaction with Azure Resource Manager `postgresql` (API Version `2023-06-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2023-06-01-preview/virtualendpoints"
```


### Client Initialization

```go
client := virtualendpoints.NewVirtualEndpointsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `VirtualEndpointsClient.Create`

```go
ctx := context.TODO()
id := virtualendpoints.NewVirtualEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerName", "virtualEndpointName")

payload := virtualendpoints.VirtualEndpointResource{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualEndpointsClient.Delete`

```go
ctx := context.TODO()
id := virtualendpoints.NewVirtualEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerName", "virtualEndpointName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualEndpointsClient.Get`

```go
ctx := context.TODO()
id := virtualendpoints.NewVirtualEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerName", "virtualEndpointName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualEndpointsClient.ListByServer`

```go
ctx := context.TODO()
id := virtualendpoints.NewFlexibleServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerName")

// alternatively `client.ListByServer(ctx, id)` can be used to do batched pagination
items, err := client.ListByServerComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualEndpointsClient.Update`

```go
ctx := context.TODO()
id := virtualendpoints.NewVirtualEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerName", "virtualEndpointName")

payload := virtualendpoints.VirtualEndpointResourceForPatch{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

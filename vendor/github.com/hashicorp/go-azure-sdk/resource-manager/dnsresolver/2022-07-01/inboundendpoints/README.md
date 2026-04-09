
## `github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/inboundendpoints` Documentation

The `inboundendpoints` SDK allows for interaction with Azure Resource Manager `dnsresolver` (API Version `2022-07-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/inboundendpoints"
```


### Client Initialization

```go
client := inboundendpoints.NewInboundEndpointsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `InboundEndpointsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := inboundendpoints.NewInboundEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dnsResolverName", "inboundEndpointName")

payload := inboundendpoints.InboundEndpoint{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload, inboundendpoints.DefaultCreateOrUpdateOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `InboundEndpointsClient.Delete`

```go
ctx := context.TODO()
id := inboundendpoints.NewInboundEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dnsResolverName", "inboundEndpointName")

if err := client.DeleteThenPoll(ctx, id, inboundendpoints.DefaultDeleteOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `InboundEndpointsClient.Get`

```go
ctx := context.TODO()
id := inboundendpoints.NewInboundEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dnsResolverName", "inboundEndpointName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `InboundEndpointsClient.List`

```go
ctx := context.TODO()
id := inboundendpoints.NewDnsResolverID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dnsResolverName")

// alternatively `client.List(ctx, id, inboundendpoints.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, inboundendpoints.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `InboundEndpointsClient.Update`

```go
ctx := context.TODO()
id := inboundendpoints.NewInboundEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dnsResolverName", "inboundEndpointName")

payload := inboundendpoints.InboundEndpointPatch{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload, inboundendpoints.DefaultUpdateOperationOptions()); err != nil {
	// handle the error
}
```

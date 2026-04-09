
## `github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/outboundendpoints` Documentation

The `outboundendpoints` SDK allows for interaction with Azure Resource Manager `dnsresolver` (API Version `2022-07-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/outboundendpoints"
```


### Client Initialization

```go
client := outboundendpoints.NewOutboundEndpointsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `OutboundEndpointsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := outboundendpoints.NewOutboundEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dnsResolverName", "outboundEndpointName")

payload := outboundendpoints.OutboundEndpoint{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload, outboundendpoints.DefaultCreateOrUpdateOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `OutboundEndpointsClient.Delete`

```go
ctx := context.TODO()
id := outboundendpoints.NewOutboundEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dnsResolverName", "outboundEndpointName")

if err := client.DeleteThenPoll(ctx, id, outboundendpoints.DefaultDeleteOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `OutboundEndpointsClient.Get`

```go
ctx := context.TODO()
id := outboundendpoints.NewOutboundEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dnsResolverName", "outboundEndpointName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OutboundEndpointsClient.List`

```go
ctx := context.TODO()
id := outboundendpoints.NewDnsResolverID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dnsResolverName")

// alternatively `client.List(ctx, id, outboundendpoints.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, outboundendpoints.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OutboundEndpointsClient.Update`

```go
ctx := context.TODO()
id := outboundendpoints.NewOutboundEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dnsResolverName", "outboundEndpointName")

payload := outboundendpoints.OutboundEndpointPatch{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload, outboundendpoints.DefaultUpdateOperationOptions()); err != nil {
	// handle the error
}
```

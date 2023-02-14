
## `github.com/hashicorp/go-azure-sdk/resource-manager/media/2022-08-01/streamingendpoints` Documentation

The `streamingendpoints` SDK allows for interaction with the Azure Resource Manager Service `media` (API Version `2022-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/media/2022-08-01/streamingendpoints"
```


### Client Initialization

```go
client := streamingendpoints.NewStreamingEndpointsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `StreamingEndpointsClient.AsyncOperation`

```go
ctx := context.TODO()
id := streamingendpoints.NewStreamingEndpointOperationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "operationIdValue")

read, err := client.AsyncOperation(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StreamingEndpointsClient.Create`

```go
ctx := context.TODO()
id := streamingendpoints.NewStreamingEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "streamingEndpointValue")

payload := streamingendpoints.StreamingEndpoint{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload, streamingendpoints.DefaultCreateOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `StreamingEndpointsClient.Delete`

```go
ctx := context.TODO()
id := streamingendpoints.NewStreamingEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "streamingEndpointValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `StreamingEndpointsClient.Get`

```go
ctx := context.TODO()
id := streamingendpoints.NewStreamingEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "streamingEndpointValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StreamingEndpointsClient.List`

```go
ctx := context.TODO()
id := streamingendpoints.NewMediaServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `StreamingEndpointsClient.OperationLocation`

```go
ctx := context.TODO()
id := streamingendpoints.NewStreamingEndpointOperationLocationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "streamingEndpointValue", "operationIdValue")

read, err := client.OperationLocation(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StreamingEndpointsClient.Scale`

```go
ctx := context.TODO()
id := streamingendpoints.NewStreamingEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "streamingEndpointValue")

payload := streamingendpoints.StreamingEntityScaleUnit{
	// ...
}


if err := client.ScaleThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `StreamingEndpointsClient.Skus`

```go
ctx := context.TODO()
id := streamingendpoints.NewStreamingEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "streamingEndpointValue")

read, err := client.Skus(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StreamingEndpointsClient.Start`

```go
ctx := context.TODO()
id := streamingendpoints.NewStreamingEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "streamingEndpointValue")

if err := client.StartThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `StreamingEndpointsClient.Stop`

```go
ctx := context.TODO()
id := streamingendpoints.NewStreamingEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "streamingEndpointValue")

if err := client.StopThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `StreamingEndpointsClient.Update`

```go
ctx := context.TODO()
id := streamingendpoints.NewStreamingEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "streamingEndpointValue")

payload := streamingendpoints.StreamingEndpoint{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

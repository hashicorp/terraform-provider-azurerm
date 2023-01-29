
## `github.com/hashicorp/go-azure-sdk/resource-manager/media/2022-08-01/streamingendpoint` Documentation

The `streamingendpoint` SDK allows for interaction with the Azure Resource Manager Service `media` (API Version `2022-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/media/2022-08-01/streamingendpoint"
```


### Client Initialization

```go
client := streamingendpoint.NewStreamingEndpointClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `StreamingEndpointClient.Update`

```go
ctx := context.TODO()
id := streamingendpoint.NewStreamingEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mediaServiceValue", "streamingEndpointValue")

payload := streamingendpoint.StreamingEndpoint{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

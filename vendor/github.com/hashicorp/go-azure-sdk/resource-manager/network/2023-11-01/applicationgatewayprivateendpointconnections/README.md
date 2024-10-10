
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/applicationgatewayprivateendpointconnections` Documentation

The `applicationgatewayprivateendpointconnections` SDK allows for interaction with Azure Resource Manager `network` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/applicationgatewayprivateendpointconnections"
```


### Client Initialization

```go
client := applicationgatewayprivateendpointconnections.NewApplicationGatewayPrivateEndpointConnectionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ApplicationGatewayPrivateEndpointConnectionsClient.Delete`

```go
ctx := context.TODO()
id := applicationgatewayprivateendpointconnections.NewApplicationGatewayPrivateEndpointConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "applicationGatewayName", "privateEndpointConnectionName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ApplicationGatewayPrivateEndpointConnectionsClient.Get`

```go
ctx := context.TODO()
id := applicationgatewayprivateendpointconnections.NewApplicationGatewayPrivateEndpointConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "applicationGatewayName", "privateEndpointConnectionName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApplicationGatewayPrivateEndpointConnectionsClient.List`

```go
ctx := context.TODO()
id := applicationgatewayprivateendpointconnections.NewApplicationGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "applicationGatewayName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ApplicationGatewayPrivateEndpointConnectionsClient.Update`

```go
ctx := context.TODO()
id := applicationgatewayprivateendpointconnections.NewApplicationGatewayPrivateEndpointConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "applicationGatewayName", "privateEndpointConnectionName")

payload := applicationgatewayprivateendpointconnections.ApplicationGatewayPrivateEndpointConnection{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/applicationgatewayprivateendpointconnections` Documentation

The `applicationgatewayprivateendpointconnections` SDK allows for interaction with the Azure Resource Manager Service `network` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/applicationgatewayprivateendpointconnections"
```


### Client Initialization

```go
client := applicationgatewayprivateendpointconnections.NewApplicationGatewayPrivateEndpointConnectionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ApplicationGatewayPrivateEndpointConnectionsClient.Delete`

```go
ctx := context.TODO()
id := applicationgatewayprivateendpointconnections.NewApplicationGatewayPrivateEndpointConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "applicationGatewayValue", "privateEndpointConnectionValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ApplicationGatewayPrivateEndpointConnectionsClient.Get`

```go
ctx := context.TODO()
id := applicationgatewayprivateendpointconnections.NewApplicationGatewayPrivateEndpointConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "applicationGatewayValue", "privateEndpointConnectionValue")

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
id := applicationgatewayprivateendpointconnections.NewApplicationGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "applicationGatewayValue")

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
id := applicationgatewayprivateendpointconnections.NewApplicationGatewayPrivateEndpointConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "applicationGatewayValue", "privateEndpointConnectionValue")

payload := applicationgatewayprivateendpointconnections.ApplicationGatewayPrivateEndpointConnection{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

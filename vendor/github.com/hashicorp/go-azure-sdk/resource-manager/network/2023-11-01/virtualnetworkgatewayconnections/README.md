
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/virtualnetworkgatewayconnections` Documentation

The `virtualnetworkgatewayconnections` SDK allows for interaction with Azure Resource Manager `network` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/virtualnetworkgatewayconnections"
```


### Client Initialization

```go
client := virtualnetworkgatewayconnections.NewVirtualNetworkGatewayConnectionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `VirtualNetworkGatewayConnectionsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := virtualnetworkgatewayconnections.NewConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "connectionName")

payload := virtualnetworkgatewayconnections.VirtualNetworkGatewayConnection{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualNetworkGatewayConnectionsClient.Delete`

```go
ctx := context.TODO()
id := virtualnetworkgatewayconnections.NewConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "connectionName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualNetworkGatewayConnectionsClient.Get`

```go
ctx := context.TODO()
id := virtualnetworkgatewayconnections.NewConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "connectionName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualNetworkGatewayConnectionsClient.GetIkeSas`

```go
ctx := context.TODO()
id := virtualnetworkgatewayconnections.NewConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "connectionName")

if err := client.GetIkeSasThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualNetworkGatewayConnectionsClient.GetSharedKey`

```go
ctx := context.TODO()
id := virtualnetworkgatewayconnections.NewConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "connectionName")

read, err := client.GetSharedKey(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualNetworkGatewayConnectionsClient.List`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VirtualNetworkGatewayConnectionsClient.ResetConnection`

```go
ctx := context.TODO()
id := virtualnetworkgatewayconnections.NewConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "connectionName")

if err := client.ResetConnectionThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualNetworkGatewayConnectionsClient.ResetSharedKey`

```go
ctx := context.TODO()
id := virtualnetworkgatewayconnections.NewConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "connectionName")

payload := virtualnetworkgatewayconnections.ConnectionResetSharedKey{
	// ...
}


if err := client.ResetSharedKeyThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualNetworkGatewayConnectionsClient.SetSharedKey`

```go
ctx := context.TODO()
id := virtualnetworkgatewayconnections.NewConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "connectionName")

payload := virtualnetworkgatewayconnections.ConnectionSharedKey{
	// ...
}


if err := client.SetSharedKeyThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualNetworkGatewayConnectionsClient.StartPacketCapture`

```go
ctx := context.TODO()
id := virtualnetworkgatewayconnections.NewConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "connectionName")

payload := virtualnetworkgatewayconnections.VpnPacketCaptureStartParameters{
	// ...
}


if err := client.StartPacketCaptureThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualNetworkGatewayConnectionsClient.StopPacketCapture`

```go
ctx := context.TODO()
id := virtualnetworkgatewayconnections.NewConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "connectionName")

payload := virtualnetworkgatewayconnections.VpnPacketCaptureStopParameters{
	// ...
}


if err := client.StopPacketCaptureThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualNetworkGatewayConnectionsClient.UpdateTags`

```go
ctx := context.TODO()
id := virtualnetworkgatewayconnections.NewConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "connectionName")

payload := virtualnetworkgatewayconnections.TagsObject{
	// ...
}


if err := client.UpdateTagsThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/vpngateways` Documentation

The `vpngateways` SDK allows for interaction with Azure Resource Manager `network` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/vpngateways"
```


### Client Initialization

```go
client := vpngateways.NewVpnGatewaysClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `VpnGatewaysClient.Reset`

```go
ctx := context.TODO()
id := vpngateways.NewVpnGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnGatewayName")

if err := client.ResetThenPoll(ctx, id, vpngateways.DefaultResetOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `VpnGatewaysClient.StartPacketCapture`

```go
ctx := context.TODO()
id := vpngateways.NewVpnGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnGatewayName")

payload := vpngateways.VpnGatewayPacketCaptureStartParameters{
	// ...
}


if err := client.StartPacketCaptureThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VpnGatewaysClient.StopPacketCapture`

```go
ctx := context.TODO()
id := vpngateways.NewVpnGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnGatewayName")

payload := vpngateways.VpnGatewayPacketCaptureStopParameters{
	// ...
}


if err := client.StopPacketCaptureThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VpnGatewaysClient.UpdateTags`

```go
ctx := context.TODO()
id := vpngateways.NewVpnGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnGatewayName")

payload := vpngateways.TagsObject{
	// ...
}


if err := client.UpdateTagsThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

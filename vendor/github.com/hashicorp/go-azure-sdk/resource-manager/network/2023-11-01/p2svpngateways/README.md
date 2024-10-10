
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/p2svpngateways` Documentation

The `p2svpngateways` SDK allows for interaction with Azure Resource Manager `network` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/p2svpngateways"
```


### Client Initialization

```go
client := p2svpngateways.NewP2sVpnGatewaysClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `P2sVpnGatewaysClient.DisconnectP2sVpnConnections`

```go
ctx := context.TODO()
id := commonids.NewVirtualWANP2SVPNGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "p2sVpnGatewayName")

payload := p2svpngateways.P2SVpnConnectionRequest{
	// ...
}


if err := client.DisconnectP2sVpnConnectionsThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `P2sVpnGatewaysClient.GenerateVpnProfile`

```go
ctx := context.TODO()
id := commonids.NewVirtualWANP2SVPNGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "p2sVpnGatewayName")

payload := p2svpngateways.P2SVpnProfileParameters{
	// ...
}


if err := client.GenerateVpnProfileThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `P2sVpnGatewaysClient.GetP2sVpnConnectionHealth`

```go
ctx := context.TODO()
id := commonids.NewVirtualWANP2SVPNGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "p2sVpnGatewayName")

if err := client.GetP2sVpnConnectionHealthThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `P2sVpnGatewaysClient.GetP2sVpnConnectionHealthDetailed`

```go
ctx := context.TODO()
id := commonids.NewVirtualWANP2SVPNGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "p2sVpnGatewayName")

payload := p2svpngateways.P2SVpnConnectionHealthRequest{
	// ...
}


if err := client.GetP2sVpnConnectionHealthDetailedThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `P2sVpnGatewaysClient.Reset`

```go
ctx := context.TODO()
id := commonids.NewVirtualWANP2SVPNGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "p2sVpnGatewayName")

if err := client.ResetThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `P2sVpnGatewaysClient.UpdateTags`

```go
ctx := context.TODO()
id := commonids.NewVirtualWANP2SVPNGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "p2sVpnGatewayName")

payload := p2svpngateways.TagsObject{
	// ...
}


if err := client.UpdateTagsThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

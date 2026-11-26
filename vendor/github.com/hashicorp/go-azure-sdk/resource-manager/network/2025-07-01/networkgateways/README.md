
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkgateways` Documentation

The `networkgateways` SDK allows for interaction with Azure Resource Manager `network` (API Version `2025-07-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkgateways"
```


### Client Initialization

```go
client := networkgateways.NewNetworkGatewaysClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `NetworkGatewaysClient.VirtualNetworkGatewayNatRulesCreateOrUpdate`

```go
ctx := context.TODO()
id := networkgateways.NewVirtualNetworkGatewayNatRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualNetworkGatewayName", "natRuleName")

payload := networkgateways.VirtualNetworkGatewayNatRule{
	// ...
}


if err := client.VirtualNetworkGatewayNatRulesCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `NetworkGatewaysClient.VirtualNetworkGatewayNatRulesDelete`

```go
ctx := context.TODO()
id := networkgateways.NewVirtualNetworkGatewayNatRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualNetworkGatewayName", "natRuleName")

if err := client.VirtualNetworkGatewayNatRulesDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `NetworkGatewaysClient.VirtualNetworkGatewayNatRulesGet`

```go
ctx := context.TODO()
id := networkgateways.NewVirtualNetworkGatewayNatRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualNetworkGatewayName", "natRuleName")

read, err := client.VirtualNetworkGatewayNatRulesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NetworkGatewaysClient.VirtualNetworkGatewayNatRulesListByVirtualNetworkGateway`

```go
ctx := context.TODO()
id := networkgateways.NewVirtualNetworkGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualNetworkGatewayName")

// alternatively `client.VirtualNetworkGatewayNatRulesListByVirtualNetworkGateway(ctx, id)` can be used to do batched pagination
items, err := client.VirtualNetworkGatewayNatRulesListByVirtualNetworkGatewayComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NetworkGatewaysClient.VirtualNetworkGatewaysInvokeAbortMigration`

```go
ctx := context.TODO()
id := networkgateways.NewVirtualNetworkGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualNetworkGatewayName")

if err := client.VirtualNetworkGatewaysInvokeAbortMigrationThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `NetworkGatewaysClient.VirtualNetworkGatewaysInvokeCommitMigration`

```go
ctx := context.TODO()
id := networkgateways.NewVirtualNetworkGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualNetworkGatewayName")

if err := client.VirtualNetworkGatewaysInvokeCommitMigrationThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `NetworkGatewaysClient.VirtualNetworkGatewaysInvokeExecuteMigration`

```go
ctx := context.TODO()
id := networkgateways.NewVirtualNetworkGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualNetworkGatewayName")

if err := client.VirtualNetworkGatewaysInvokeExecuteMigrationThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `NetworkGatewaysClient.VirtualNetworkGatewaysInvokePrepareMigration`

```go
ctx := context.TODO()
id := networkgateways.NewVirtualNetworkGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualNetworkGatewayName")

payload := networkgateways.VirtualNetworkGatewayMigrationParameters{
	// ...
}


if err := client.VirtualNetworkGatewaysInvokePrepareMigrationThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

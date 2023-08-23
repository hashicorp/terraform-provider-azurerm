
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/localnetworkgateways` Documentation

The `localnetworkgateways` SDK allows for interaction with the Azure Resource Manager Service `network` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/localnetworkgateways"
```


### Client Initialization

```go
client := localnetworkgateways.NewLocalNetworkGatewaysClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `LocalNetworkGatewaysClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := localnetworkgateways.NewLocalNetworkGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localNetworkGatewayValue")

payload := localnetworkgateways.LocalNetworkGateway{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `LocalNetworkGatewaysClient.Delete`

```go
ctx := context.TODO()
id := localnetworkgateways.NewLocalNetworkGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localNetworkGatewayValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `LocalNetworkGatewaysClient.Get`

```go
ctx := context.TODO()
id := localnetworkgateways.NewLocalNetworkGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localNetworkGatewayValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LocalNetworkGatewaysClient.List`

```go
ctx := context.TODO()
id := localnetworkgateways.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `LocalNetworkGatewaysClient.UpdateTags`

```go
ctx := context.TODO()
id := localnetworkgateways.NewLocalNetworkGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "localNetworkGatewayValue")

payload := localnetworkgateways.TagsObject{
	// ...
}


read, err := client.UpdateTags(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/applicationgatewayprivatelinkresources` Documentation

The `applicationgatewayprivatelinkresources` SDK allows for interaction with the Azure Resource Manager Service `network` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/applicationgatewayprivatelinkresources"
```


### Client Initialization

```go
client := applicationgatewayprivatelinkresources.NewApplicationGatewayPrivateLinkResourcesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ApplicationGatewayPrivateLinkResourcesClient.List`

```go
ctx := context.TODO()
id := applicationgatewayprivatelinkresources.NewApplicationGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "applicationGatewayValue")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

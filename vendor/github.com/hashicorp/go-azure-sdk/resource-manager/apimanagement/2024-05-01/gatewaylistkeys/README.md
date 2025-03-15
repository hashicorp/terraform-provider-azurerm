
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/gatewaylistkeys` Documentation

The `gatewaylistkeys` SDK allows for interaction with Azure Resource Manager `apimanagement` (API Version `2024-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/gatewaylistkeys"
```


### Client Initialization

```go
client := gatewaylistkeys.NewGatewayListKeysClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `GatewayListKeysClient.GatewayListKeys`

```go
ctx := context.TODO()
id := gatewaylistkeys.NewServiceGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "gatewayId")

read, err := client.GatewayListKeys(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

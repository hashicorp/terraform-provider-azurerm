
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/gateway` Documentation

The `gateway` SDK allows for interaction with the Azure Resource Manager Service `apimanagement` (API Version `2021-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/gateway"
```


### Client Initialization

```go
client := gateway.NewGatewayClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `GatewayClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := gateway.NewGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "gatewayIdValue")

payload := gateway.GatewayContract{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, gateway.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GatewayClient.Delete`

```go
ctx := context.TODO()
id := gateway.NewGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "gatewayIdValue")

read, err := client.Delete(ctx, id, gateway.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GatewayClient.Get`

```go
ctx := context.TODO()
id := gateway.NewGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "gatewayIdValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GatewayClient.GetEntityTag`

```go
ctx := context.TODO()
id := gateway.NewGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "gatewayIdValue")

read, err := client.GetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GatewayClient.ListByService`

```go
ctx := context.TODO()
id := gateway.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue")

// alternatively `client.ListByService(ctx, id, gateway.DefaultListByServiceOperationOptions())` can be used to do batched pagination
items, err := client.ListByServiceComplete(ctx, id, gateway.DefaultListByServiceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `GatewayClient.Update`

```go
ctx := context.TODO()
id := gateway.NewGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "gatewayIdValue")

payload := gateway.GatewayContract{
	// ...
}


read, err := client.Update(ctx, id, payload, gateway.DefaultUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

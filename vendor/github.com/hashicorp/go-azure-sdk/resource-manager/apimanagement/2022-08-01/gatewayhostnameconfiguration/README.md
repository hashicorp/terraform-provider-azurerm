
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/gatewayhostnameconfiguration` Documentation

The `gatewayhostnameconfiguration` SDK allows for interaction with Azure Resource Manager `apimanagement` (API Version `2022-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/gatewayhostnameconfiguration"
```


### Client Initialization

```go
client := gatewayhostnameconfiguration.NewGatewayHostnameConfigurationClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `GatewayHostnameConfigurationClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := gatewayhostnameconfiguration.NewHostnameConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "gatewayId", "hcId")

payload := gatewayhostnameconfiguration.GatewayHostnameConfigurationContract{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, gatewayhostnameconfiguration.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GatewayHostnameConfigurationClient.Delete`

```go
ctx := context.TODO()
id := gatewayhostnameconfiguration.NewHostnameConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "gatewayId", "hcId")

read, err := client.Delete(ctx, id, gatewayhostnameconfiguration.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GatewayHostnameConfigurationClient.Get`

```go
ctx := context.TODO()
id := gatewayhostnameconfiguration.NewHostnameConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "gatewayId", "hcId")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GatewayHostnameConfigurationClient.GetEntityTag`

```go
ctx := context.TODO()
id := gatewayhostnameconfiguration.NewHostnameConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "gatewayId", "hcId")

read, err := client.GetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GatewayHostnameConfigurationClient.ListByService`

```go
ctx := context.TODO()
id := gatewayhostnameconfiguration.NewGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "gatewayId")

// alternatively `client.ListByService(ctx, id, gatewayhostnameconfiguration.DefaultListByServiceOperationOptions())` can be used to do batched pagination
items, err := client.ListByServiceComplete(ctx, id, gatewayhostnameconfiguration.DefaultListByServiceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

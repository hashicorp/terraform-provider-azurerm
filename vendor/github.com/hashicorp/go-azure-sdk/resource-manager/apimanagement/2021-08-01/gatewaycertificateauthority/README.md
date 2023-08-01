
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/gatewaycertificateauthority` Documentation

The `gatewaycertificateauthority` SDK allows for interaction with the Azure Resource Manager Service `apimanagement` (API Version `2021-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/gatewaycertificateauthority"
```


### Client Initialization

```go
client := gatewaycertificateauthority.NewGatewayCertificateAuthorityClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `GatewayCertificateAuthorityClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := gatewaycertificateauthority.NewCertificateAuthorityID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "gatewayIdValue", "certificateIdValue")

payload := gatewaycertificateauthority.GatewayCertificateAuthorityContract{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, gatewaycertificateauthority.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GatewayCertificateAuthorityClient.Delete`

```go
ctx := context.TODO()
id := gatewaycertificateauthority.NewCertificateAuthorityID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "gatewayIdValue", "certificateIdValue")

read, err := client.Delete(ctx, id, gatewaycertificateauthority.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GatewayCertificateAuthorityClient.Get`

```go
ctx := context.TODO()
id := gatewaycertificateauthority.NewCertificateAuthorityID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "gatewayIdValue", "certificateIdValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GatewayCertificateAuthorityClient.GetEntityTag`

```go
ctx := context.TODO()
id := gatewaycertificateauthority.NewCertificateAuthorityID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "gatewayIdValue", "certificateIdValue")

read, err := client.GetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GatewayCertificateAuthorityClient.ListByService`

```go
ctx := context.TODO()
id := gatewaycertificateauthority.NewGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceValue", "gatewayIdValue")

// alternatively `client.ListByService(ctx, id, gatewaycertificateauthority.DefaultListByServiceOperationOptions())` can be used to do batched pagination
items, err := client.ListByServiceComplete(ctx, id, gatewaycertificateauthority.DefaultListByServiceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

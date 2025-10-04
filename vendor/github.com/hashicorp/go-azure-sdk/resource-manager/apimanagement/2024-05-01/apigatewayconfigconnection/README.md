
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apigatewayconfigconnection` Documentation

The `apigatewayconfigconnection` SDK allows for interaction with Azure Resource Manager `apimanagement` (API Version `2024-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apigatewayconfigconnection"
```


### Client Initialization

```go
client := apigatewayconfigconnection.NewApiGatewayConfigConnectionClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ApiGatewayConfigConnectionClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := apigatewayconfigconnection.NewConfigConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "gatewayName", "configConnectionName")

payload := apigatewayconfigconnection.ApiManagementGatewayConfigConnectionResource{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ApiGatewayConfigConnectionClient.Delete`

```go
ctx := context.TODO()
id := apigatewayconfigconnection.NewConfigConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "gatewayName", "configConnectionName")

if err := client.DeleteThenPoll(ctx, id, apigatewayconfigconnection.DefaultDeleteOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `ApiGatewayConfigConnectionClient.Get`

```go
ctx := context.TODO()
id := apigatewayconfigconnection.NewConfigConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "gatewayName", "configConnectionName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiGatewayConfigConnectionClient.ListByGateway`

```go
ctx := context.TODO()
id := apigatewayconfigconnection.NewGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "gatewayName")

// alternatively `client.ListByGateway(ctx, id)` can be used to do batched pagination
items, err := client.ListByGatewayComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

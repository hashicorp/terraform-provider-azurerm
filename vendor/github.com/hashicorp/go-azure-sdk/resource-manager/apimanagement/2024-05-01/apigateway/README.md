
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apigateway` Documentation

The `apigateway` SDK allows for interaction with Azure Resource Manager `apimanagement` (API Version `2024-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apigateway"
```


### Client Initialization

```go
client := apigateway.NewApiGatewayClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ApiGatewayClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := apigateway.NewGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "gatewayName")

payload := apigateway.ApiManagementGatewayResource{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ApiGatewayClient.Delete`

```go
ctx := context.TODO()
id := apigateway.NewGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "gatewayName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ApiGatewayClient.Get`

```go
ctx := context.TODO()
id := apigateway.NewGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "gatewayName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiGatewayClient.List`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ApiGatewayClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ApiGatewayClient.Update`

```go
ctx := context.TODO()
id := apigateway.NewGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "gatewayName")

payload := apigateway.ApiManagementGatewayUpdateParameters{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

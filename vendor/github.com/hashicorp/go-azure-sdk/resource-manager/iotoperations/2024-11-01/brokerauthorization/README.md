
## `github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/brokerauthorization` Documentation

The `brokerauthorization` SDK allows for interaction with Azure Resource Manager `iotoperations` (API Version `2024-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/brokerauthorization"
```


### Client Initialization

```go
client := brokerauthorization.NewBrokerAuthorizationClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `BrokerAuthorizationClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := brokerauthorization.NewAuthorizationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "instanceName", "brokerName", "authorizationName")

payload := brokerauthorization.BrokerAuthorizationResource{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `BrokerAuthorizationClient.Delete`

```go
ctx := context.TODO()
id := brokerauthorization.NewAuthorizationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "instanceName", "brokerName", "authorizationName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `BrokerAuthorizationClient.Get`

```go
ctx := context.TODO()
id := brokerauthorization.NewAuthorizationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "instanceName", "brokerName", "authorizationName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BrokerAuthorizationClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := brokerauthorization.NewBrokerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "instanceName", "brokerName")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


## `github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/broker` Documentation

The `broker` SDK allows for interaction with Azure Resource Manager `iotoperations` (API Version `2024-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/broker"
```


### Client Initialization

```go
client := broker.NewBrokerClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `BrokerClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := broker.NewBrokerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "instanceName", "brokerName")

payload := broker.BrokerResource{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `BrokerClient.Delete`

```go
ctx := context.TODO()
id := broker.NewBrokerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "instanceName", "brokerName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `BrokerClient.Get`

```go
ctx := context.TODO()
id := broker.NewBrokerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "instanceName", "brokerName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BrokerClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := broker.NewInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "instanceName")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

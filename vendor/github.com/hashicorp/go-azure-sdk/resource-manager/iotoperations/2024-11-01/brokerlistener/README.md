
## `github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/brokerlistener` Documentation

The `brokerlistener` SDK allows for interaction with Azure Resource Manager `iotoperations` (API Version `2024-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/brokerlistener"
```


### Client Initialization

```go
client := brokerlistener.NewBrokerListenerClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `BrokerListenerClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := brokerlistener.NewListenerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "instanceName", "brokerName", "listenerName")

payload := brokerlistener.BrokerListenerResource{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `BrokerListenerClient.Delete`

```go
ctx := context.TODO()
id := brokerlistener.NewListenerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "instanceName", "brokerName", "listenerName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `BrokerListenerClient.Get`

```go
ctx := context.TODO()
id := brokerlistener.NewListenerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "instanceName", "brokerName", "listenerName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BrokerListenerClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := brokerlistener.NewBrokerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "instanceName", "brokerName")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


## `github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/consumergroups` Documentation

The `consumergroups` SDK allows for interaction with Azure Resource Manager `eventhub` (API Version `2024-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/consumergroups"
```


### Client Initialization

```go
client := consumergroups.NewConsumerGroupsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ConsumerGroupsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := consumergroups.NewConsumerGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "eventhubName", "consumerGroupName")

payload := consumergroups.ConsumerGroup{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConsumerGroupsClient.Delete`

```go
ctx := context.TODO()
id := consumergroups.NewConsumerGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "eventhubName", "consumerGroupName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConsumerGroupsClient.Get`

```go
ctx := context.TODO()
id := consumergroups.NewConsumerGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "eventhubName", "consumerGroupName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConsumerGroupsClient.ListByEventHub`

```go
ctx := context.TODO()
id := consumergroups.NewEventhubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "eventhubName")

// alternatively `client.ListByEventHub(ctx, id, consumergroups.DefaultListByEventHubOperationOptions())` can be used to do batched pagination
items, err := client.ListByEventHubComplete(ctx, id, consumergroups.DefaultListByEventHubOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


## `github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/queues` Documentation

The `queues` SDK allows for interaction with Azure Resource Manager `servicebus` (API Version `2021-06-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/queues"
```


### Client Initialization

```go
client := queues.NewQueuesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `QueuesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := queues.NewQueueID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "queueName")

payload := queues.SBQueue{
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


### Example Usage: `QueuesClient.Delete`

```go
ctx := context.TODO()
id := queues.NewQueueID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "queueName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `QueuesClient.Get`

```go
ctx := context.TODO()
id := queues.NewQueueID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "queueName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `QueuesClient.ListByNamespace`

```go
ctx := context.TODO()
id := queues.NewNamespaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName")

// alternatively `client.ListByNamespace(ctx, id, queues.DefaultListByNamespaceOperationOptions())` can be used to do batched pagination
items, err := client.ListByNamespaceComplete(ctx, id, queues.DefaultListByNamespaceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

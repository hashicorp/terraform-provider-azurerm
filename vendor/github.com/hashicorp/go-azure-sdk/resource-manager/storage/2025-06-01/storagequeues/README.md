
## `github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-06-01/storagequeues` Documentation

The `storagequeues` SDK allows for interaction with Azure Resource Manager `storage` (API Version `2025-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-06-01/storagequeues"
```


### Client Initialization

```go
client := storagequeues.NewStorageQueuesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `StorageQueuesClient.QueueCreate`

```go
ctx := context.TODO()
id := storagequeues.NewQueueID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName", "queueName")

payload := storagequeues.StorageQueue{
	// ...
}


read, err := client.QueueCreate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StorageQueuesClient.QueueDelete`

```go
ctx := context.TODO()
id := storagequeues.NewQueueID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName", "queueName")

read, err := client.QueueDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StorageQueuesClient.QueueGet`

```go
ctx := context.TODO()
id := storagequeues.NewQueueID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName", "queueName")

read, err := client.QueueGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StorageQueuesClient.QueueUpdate`

```go
ctx := context.TODO()
id := storagequeues.NewQueueID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName", "queueName")

payload := storagequeues.StorageQueue{
	// ...
}


read, err := client.QueueUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

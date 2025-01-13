
## `github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/queueservice` Documentation

The `queueservice` SDK allows for interaction with Azure Resource Manager `storage` (API Version `2023-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/queueservice"
```


### Client Initialization

```go
client := queueservice.NewQueueServiceClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `QueueServiceClient.QueueCreate`

```go
ctx := context.TODO()
id := queueservice.NewQueueID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName", "queueName")

payload := queueservice.StorageQueue{
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


### Example Usage: `QueueServiceClient.QueueDelete`

```go
ctx := context.TODO()
id := queueservice.NewQueueID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName", "queueName")

read, err := client.QueueDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `QueueServiceClient.QueueGet`

```go
ctx := context.TODO()
id := queueservice.NewQueueID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName", "queueName")

read, err := client.QueueGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `QueueServiceClient.QueueList`

```go
ctx := context.TODO()
id := commonids.NewStorageAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName")

// alternatively `client.QueueList(ctx, id, queueservice.DefaultQueueListOperationOptions())` can be used to do batched pagination
items, err := client.QueueListComplete(ctx, id, queueservice.DefaultQueueListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `QueueServiceClient.QueueUpdate`

```go
ctx := context.TODO()
id := queueservice.NewQueueID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName", "queueName")

payload := queueservice.StorageQueue{
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

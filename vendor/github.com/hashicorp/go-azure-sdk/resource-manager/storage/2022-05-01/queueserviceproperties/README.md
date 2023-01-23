
## `github.com/hashicorp/go-azure-sdk/resource-manager/storage/2022-05-01/queueserviceproperties` Documentation

The `queueserviceproperties` SDK allows for interaction with the Azure Resource Manager Service `storage` (API Version `2022-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/storage/2022-05-01/queueserviceproperties"
```


### Client Initialization

```go
client := queueserviceproperties.NewQueueServicePropertiesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `QueueServicePropertiesClient.QueueServicesGetServiceProperties`

```go
ctx := context.TODO()
id := queueserviceproperties.NewStorageAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountValue")

read, err := client.QueueServicesGetServiceProperties(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `QueueServicePropertiesClient.QueueServicesList`

```go
ctx := context.TODO()
id := queueserviceproperties.NewStorageAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountValue")

read, err := client.QueueServicesList(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `QueueServicePropertiesClient.QueueServicesSetServiceProperties`

```go
ctx := context.TODO()
id := queueserviceproperties.NewStorageAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountValue")

payload := queueserviceproperties.QueueServiceProperties{
	// ...
}


read, err := client.QueueServicesSetServiceProperties(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

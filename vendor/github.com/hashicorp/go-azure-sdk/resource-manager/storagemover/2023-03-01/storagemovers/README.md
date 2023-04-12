
## `github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2023-03-01/storagemovers` Documentation

The `storagemovers` SDK allows for interaction with the Azure Resource Manager Service `storagemover` (API Version `2023-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2023-03-01/storagemovers"
```


### Client Initialization

```go
client := storagemovers.NewStorageMoversClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `StorageMoversClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := storagemovers.NewStorageMoverID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageMoverValue")

payload := storagemovers.StorageMover{
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


### Example Usage: `StorageMoversClient.Delete`

```go
ctx := context.TODO()
id := storagemovers.NewStorageMoverID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageMoverValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `StorageMoversClient.Get`

```go
ctx := context.TODO()
id := storagemovers.NewStorageMoverID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageMoverValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StorageMoversClient.List`

```go
ctx := context.TODO()
id := storagemovers.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `StorageMoversClient.ListBySubscription`

```go
ctx := context.TODO()
id := storagemovers.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `StorageMoversClient.Update`

```go
ctx := context.TODO()
id := storagemovers.NewStorageMoverID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageMoverValue")

payload := storagemovers.StorageMoverUpdateParameters{
	// ...
}


read, err := client.Update(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

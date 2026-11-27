
## `github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/storageassets` Documentation

The `storageassets` SDK allows for interaction with Azure Resource Manager `discovery` (API Version `2026-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/storageassets"
```


### Client Initialization

```go
client := storageassets.NewStorageAssetsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `StorageAssetsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := storageassets.NewStorageAssetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageContainerName", "storageAssetName")

payload := storageassets.StorageAsset{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `StorageAssetsClient.Delete`

```go
ctx := context.TODO()
id := storageassets.NewStorageAssetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageContainerName", "storageAssetName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `StorageAssetsClient.Get`

```go
ctx := context.TODO()
id := storageassets.NewStorageAssetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageContainerName", "storageAssetName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StorageAssetsClient.ListByStorageContainer`

```go
ctx := context.TODO()
id := storageassets.NewStorageContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageContainerName")

// alternatively `client.ListByStorageContainer(ctx, id)` can be used to do batched pagination
items, err := client.ListByStorageContainerComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `StorageAssetsClient.Update`

```go
ctx := context.TODO()
id := storageassets.NewStorageAssetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageContainerName", "storageAssetName")

payload := storageassets.StorageAssetUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


## `github.com/hashicorp/go-azure-sdk/data-plane/keyvault/2025-07-01/deletedstorage` Documentation

The `deletedstorage` SDK allows for interaction with <unknown source data type> `keyvault` (API Version `2025-07-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/data-plane/keyvault/2025-07-01/deletedstorage"
```


### Client Initialization

```go
client := deletedstorage.NewDeletedStorageClientWithBaseURI("")
client.Client.Authorizer = authorizer
```


### Example Usage: `DeletedStorageClient.GetDeletedSasDefinition`

```go
ctx := context.TODO()
id := deletedstorage.NewDeletedstorageSaID("deletedstorageName", "saName")

read, err := client.GetDeletedSasDefinition(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeletedStorageClient.GetDeletedSasDefinitions`

```go
ctx := context.TODO()
id := deletedstorage.NewDeletedstorageID("deletedstorageName")

// alternatively `client.GetDeletedSasDefinitions(ctx, id, deletedstorage.DefaultGetDeletedSasDefinitionsOperationOptions())` can be used to do batched pagination
items, err := client.GetDeletedSasDefinitionsComplete(ctx, id, deletedstorage.DefaultGetDeletedSasDefinitionsOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DeletedStorageClient.GetDeletedStorageAccount`

```go
ctx := context.TODO()
id := deletedstorage.NewDeletedstorageID("deletedstorageName")

read, err := client.GetDeletedStorageAccount(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeletedStorageClient.GetDeletedStorageAccounts`

```go
ctx := context.TODO()


// alternatively `client.GetDeletedStorageAccounts(ctx, deletedstorage.DefaultGetDeletedStorageAccountsOperationOptions())` can be used to do batched pagination
items, err := client.GetDeletedStorageAccountsComplete(ctx, deletedstorage.DefaultGetDeletedStorageAccountsOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DeletedStorageClient.PurgeDeletedStorageAccount`

```go
ctx := context.TODO()
id := deletedstorage.NewDeletedstorageID("deletedstorageName")

read, err := client.PurgeDeletedStorageAccount(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeletedStorageClient.RecoverDeletedSasDefinition`

```go
ctx := context.TODO()
id := deletedstorage.NewDeletedstorageSaID("deletedstorageName", "saName")

read, err := client.RecoverDeletedSasDefinition(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeletedStorageClient.RecoverDeletedStorageAccount`

```go
ctx := context.TODO()
id := deletedstorage.NewDeletedstorageID("deletedstorageName")

read, err := client.RecoverDeletedStorageAccount(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

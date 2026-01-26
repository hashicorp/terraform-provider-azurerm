
## `github.com/hashicorp/go-azure-sdk/data-plane/keyvault/2025-07-01/storage` Documentation

The `storage` SDK allows for interaction with <unknown source data type> `keyvault` (API Version `2025-07-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/data-plane/keyvault/2025-07-01/storage"
```


### Client Initialization

```go
client := storage.NewStorageClientWithBaseURI("")
client.Client.Authorizer = authorizer
```


### Example Usage: `StorageClient.BackupStorageAccount`

```go
ctx := context.TODO()
id := storage.NewStorageID("storageName")

read, err := client.BackupStorageAccount(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StorageClient.DeleteSasDefinition`

```go
ctx := context.TODO()
id := storage.NewSaID("storageName", "saName")

read, err := client.DeleteSasDefinition(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StorageClient.DeleteStorageAccount`

```go
ctx := context.TODO()
id := storage.NewStorageID("storageName")

read, err := client.DeleteStorageAccount(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StorageClient.GetSasDefinition`

```go
ctx := context.TODO()
id := storage.NewSaID("storageName", "saName")

read, err := client.GetSasDefinition(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StorageClient.GetSasDefinitions`

```go
ctx := context.TODO()
id := storage.NewStorageID("storageName")

// alternatively `client.GetSasDefinitions(ctx, id, storage.DefaultGetSasDefinitionsOperationOptions())` can be used to do batched pagination
items, err := client.GetSasDefinitionsComplete(ctx, id, storage.DefaultGetSasDefinitionsOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `StorageClient.GetStorageAccount`

```go
ctx := context.TODO()
id := storage.NewStorageID("storageName")

read, err := client.GetStorageAccount(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StorageClient.GetStorageAccounts`

```go
ctx := context.TODO()


// alternatively `client.GetStorageAccounts(ctx, storage.DefaultGetStorageAccountsOperationOptions())` can be used to do batched pagination
items, err := client.GetStorageAccountsComplete(ctx, storage.DefaultGetStorageAccountsOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `StorageClient.RegenerateStorageAccountKey`

```go
ctx := context.TODO()
id := storage.NewStorageID("storageName")

payload := storage.StorageAccountRegenerteKeyParameters{
	// ...
}


read, err := client.RegenerateStorageAccountKey(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StorageClient.RestoreStorageAccount`

```go
ctx := context.TODO()

payload := storage.StorageRestoreParameters{
	// ...
}


read, err := client.RestoreStorageAccount(ctx, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StorageClient.SetSasDefinition`

```go
ctx := context.TODO()
id := storage.NewSaID("storageName", "saName")

payload := storage.SasDefinitionCreateParameters{
	// ...
}


read, err := client.SetSasDefinition(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StorageClient.SetStorageAccount`

```go
ctx := context.TODO()
id := storage.NewStorageID("storageName")

payload := storage.StorageAccountCreateParameters{
	// ...
}


read, err := client.SetStorageAccount(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StorageClient.UpdateSasDefinition`

```go
ctx := context.TODO()
id := storage.NewSaID("storageName", "saName")

payload := storage.SasDefinitionUpdateParameters{
	// ...
}


read, err := client.UpdateSasDefinition(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StorageClient.UpdateStorageAccount`

```go
ctx := context.TODO()
id := storage.NewStorageID("storageName")

payload := storage.StorageAccountUpdateParameters{
	// ...
}


read, err := client.UpdateStorageAccount(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

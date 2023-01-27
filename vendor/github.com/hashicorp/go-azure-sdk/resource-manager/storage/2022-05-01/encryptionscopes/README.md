
## `github.com/hashicorp/go-azure-sdk/resource-manager/storage/2022-05-01/encryptionscopes` Documentation

The `encryptionscopes` SDK allows for interaction with the Azure Resource Manager Service `storage` (API Version `2022-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/storage/2022-05-01/encryptionscopes"
```


### Client Initialization

```go
client := encryptionscopes.NewEncryptionScopesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `EncryptionScopesClient.Get`

```go
ctx := context.TODO()
id := encryptionscopes.NewEncryptionScopeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountValue", "encryptionScopeValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EncryptionScopesClient.List`

```go
ctx := context.TODO()
id := encryptionscopes.NewStorageAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountValue")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `EncryptionScopesClient.Patch`

```go
ctx := context.TODO()
id := encryptionscopes.NewEncryptionScopeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountValue", "encryptionScopeValue")

payload := encryptionscopes.EncryptionScope{
	// ...
}


read, err := client.Patch(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EncryptionScopesClient.Put`

```go
ctx := context.TODO()
id := encryptionscopes.NewEncryptionScopeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountValue", "encryptionScopeValue")

payload := encryptionscopes.EncryptionScope{
	// ...
}


read, err := client.Put(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

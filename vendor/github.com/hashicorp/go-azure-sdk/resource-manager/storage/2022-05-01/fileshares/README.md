
## `github.com/hashicorp/go-azure-sdk/resource-manager/storage/2022-05-01/fileshares` Documentation

The `fileshares` SDK allows for interaction with the Azure Resource Manager Service `storage` (API Version `2022-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/storage/2022-05-01/fileshares"
```


### Client Initialization

```go
client := fileshares.NewFileSharesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `FileSharesClient.Create`

```go
ctx := context.TODO()
id := fileshares.NewShareID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountValue", "shareValue")

payload := fileshares.FileShare{
	// ...
}


read, err := client.Create(ctx, id, payload, fileshares.DefaultCreateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FileSharesClient.Delete`

```go
ctx := context.TODO()
id := fileshares.NewShareID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountValue", "shareValue")

read, err := client.Delete(ctx, id, fileshares.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FileSharesClient.Get`

```go
ctx := context.TODO()
id := fileshares.NewShareID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountValue", "shareValue")

read, err := client.Get(ctx, id, fileshares.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FileSharesClient.Lease`

```go
ctx := context.TODO()
id := fileshares.NewShareID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountValue", "shareValue")

payload := fileshares.LeaseShareRequest{
	// ...
}


read, err := client.Lease(ctx, id, payload, fileshares.DefaultLeaseOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FileSharesClient.List`

```go
ctx := context.TODO()
id := fileshares.NewStorageAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountValue")

// alternatively `client.List(ctx, id, fileshares.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, fileshares.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `FileSharesClient.Restore`

```go
ctx := context.TODO()
id := fileshares.NewShareID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountValue", "shareValue")

payload := fileshares.DeletedShare{
	// ...
}


read, err := client.Restore(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FileSharesClient.Update`

```go
ctx := context.TODO()
id := fileshares.NewShareID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountValue", "shareValue")

payload := fileshares.FileShare{
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

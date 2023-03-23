
## `github.com/hashicorp/go-azure-sdk/resource-manager/storage/2022-05-01/blobservice` Documentation

The `blobservice` SDK allows for interaction with the Azure Resource Manager Service `storage` (API Version `2022-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/storage/2022-05-01/blobservice"
```


### Client Initialization

```go
client := blobservice.NewBlobServiceClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `BlobServiceClient.GetServiceProperties`

```go
ctx := context.TODO()
id := blobservice.NewStorageAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountValue")

read, err := client.GetServiceProperties(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BlobServiceClient.List`

```go
ctx := context.TODO()
id := blobservice.NewStorageAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountValue")

read, err := client.List(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BlobServiceClient.SetServiceProperties`

```go
ctx := context.TODO()
id := blobservice.NewStorageAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountValue")

payload := blobservice.BlobServiceProperties{
	// ...
}


read, err := client.SetServiceProperties(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


## `github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-06-01/blobservices` Documentation

The `blobservices` SDK allows for interaction with Azure Resource Manager `storage` (API Version `2025-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-06-01/blobservices"
```


### Client Initialization

```go
client := blobservices.NewBlobServicesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `BlobServicesClient.BlobContainersList`

```go
ctx := context.TODO()
id := commonids.NewStorageAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName")

// alternatively `client.BlobContainersList(ctx, id, blobservices.DefaultBlobContainersListOperationOptions())` can be used to do batched pagination
items, err := client.BlobContainersListComplete(ctx, id, blobservices.DefaultBlobContainersListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `BlobServicesClient.GetServiceProperties`

```go
ctx := context.TODO()
id := commonids.NewStorageAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName")

read, err := client.GetServiceProperties(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BlobServicesClient.List`

```go
ctx := context.TODO()
id := commonids.NewStorageAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `BlobServicesClient.SetServiceProperties`

```go
ctx := context.TODO()
id := commonids.NewStorageAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName")

payload := blobservices.BlobServiceProperties{
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


## `github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-06-01/blobcontainers` Documentation

The `blobcontainers` SDK allows for interaction with Azure Resource Manager `storage` (API Version `2025-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-06-01/blobcontainers"
```


### Client Initialization

```go
client := blobcontainers.NewBlobContainersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `BlobContainersClient.ClearLegalHold`

```go
ctx := context.TODO()
id := commonids.NewStorageContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName", "containerName")

payload := blobcontainers.LegalHold{
	// ...
}


read, err := client.ClearLegalHold(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BlobContainersClient.Create`

```go
ctx := context.TODO()
id := commonids.NewStorageContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName", "containerName")

payload := blobcontainers.BlobContainer{
	// ...
}


read, err := client.Create(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BlobContainersClient.Delete`

```go
ctx := context.TODO()
id := commonids.NewStorageContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName", "containerName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BlobContainersClient.Get`

```go
ctx := context.TODO()
id := commonids.NewStorageContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName", "containerName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BlobContainersClient.Lease`

```go
ctx := context.TODO()
id := commonids.NewStorageContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName", "containerName")

payload := blobcontainers.LeaseContainerRequest{
	// ...
}


read, err := client.Lease(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BlobContainersClient.ObjectLevelWorm`

```go
ctx := context.TODO()
id := commonids.NewStorageContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName", "containerName")

if err := client.ObjectLevelWormThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `BlobContainersClient.SetLegalHold`

```go
ctx := context.TODO()
id := commonids.NewStorageContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName", "containerName")

payload := blobcontainers.LegalHold{
	// ...
}


read, err := client.SetLegalHold(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BlobContainersClient.Update`

```go
ctx := context.TODO()
id := commonids.NewStorageContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName", "containerName")

payload := blobcontainers.BlobContainer{
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

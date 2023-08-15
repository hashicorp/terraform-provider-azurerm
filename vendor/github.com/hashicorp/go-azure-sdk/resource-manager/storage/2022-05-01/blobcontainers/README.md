
## `github.com/hashicorp/go-azure-sdk/resource-manager/storage/2022-05-01/blobcontainers` Documentation

The `blobcontainers` SDK allows for interaction with the Azure Resource Manager Service `storage` (API Version `2022-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/storage/2022-05-01/blobcontainers"
```


### Client Initialization

```go
client := blobcontainers.NewBlobContainersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `BlobContainersClient.ClearLegalHold`

```go
ctx := context.TODO()
id := blobcontainers.NewStorageContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountValue", "containerValue")

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
id := blobcontainers.NewStorageContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountValue", "containerValue")

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


### Example Usage: `BlobContainersClient.CreateOrUpdateImmutabilityPolicy`

```go
ctx := context.TODO()
id := blobcontainers.NewStorageContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountValue", "containerValue")

payload := blobcontainers.ImmutabilityPolicy{
	// ...
}


read, err := client.CreateOrUpdateImmutabilityPolicy(ctx, id, payload, blobcontainers.DefaultCreateOrUpdateImmutabilityPolicyOperationOptions())
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
id := blobcontainers.NewStorageContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountValue", "containerValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BlobContainersClient.DeleteImmutabilityPolicy`

```go
ctx := context.TODO()
id := blobcontainers.NewStorageContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountValue", "containerValue")

read, err := client.DeleteImmutabilityPolicy(ctx, id, blobcontainers.DefaultDeleteImmutabilityPolicyOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BlobContainersClient.ExtendImmutabilityPolicy`

```go
ctx := context.TODO()
id := blobcontainers.NewStorageContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountValue", "containerValue")

payload := blobcontainers.ImmutabilityPolicy{
	// ...
}


read, err := client.ExtendImmutabilityPolicy(ctx, id, payload, blobcontainers.DefaultExtendImmutabilityPolicyOperationOptions())
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
id := blobcontainers.NewStorageContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountValue", "containerValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BlobContainersClient.GetImmutabilityPolicy`

```go
ctx := context.TODO()
id := blobcontainers.NewStorageContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountValue", "containerValue")

read, err := client.GetImmutabilityPolicy(ctx, id, blobcontainers.DefaultGetImmutabilityPolicyOperationOptions())
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
id := blobcontainers.NewStorageContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountValue", "containerValue")

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


### Example Usage: `BlobContainersClient.List`

```go
ctx := context.TODO()
id := blobcontainers.NewStorageAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountValue")

// alternatively `client.List(ctx, id, blobcontainers.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, blobcontainers.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `BlobContainersClient.LockImmutabilityPolicy`

```go
ctx := context.TODO()
id := blobcontainers.NewStorageContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountValue", "containerValue")

read, err := client.LockImmutabilityPolicy(ctx, id, blobcontainers.DefaultLockImmutabilityPolicyOperationOptions())
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
id := blobcontainers.NewStorageContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountValue", "containerValue")

if err := client.ObjectLevelWormThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `BlobContainersClient.SetLegalHold`

```go
ctx := context.TODO()
id := blobcontainers.NewStorageContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountValue", "containerValue")

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
id := blobcontainers.NewStorageContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountValue", "containerValue")

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

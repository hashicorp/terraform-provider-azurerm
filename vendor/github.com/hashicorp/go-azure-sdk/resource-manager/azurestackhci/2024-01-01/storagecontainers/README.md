
## `github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/storagecontainers` Documentation

The `storagecontainers` SDK allows for interaction with Azure Resource Manager `azurestackhci` (API Version `2024-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/storagecontainers"
```


### Client Initialization

```go
client := storagecontainers.NewStorageContainersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `StorageContainersClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := storagecontainers.NewStorageContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageContainerName")

payload := storagecontainers.StorageContainers{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `StorageContainersClient.Delete`

```go
ctx := context.TODO()
id := storagecontainers.NewStorageContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageContainerName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `StorageContainersClient.Get`

```go
ctx := context.TODO()
id := storagecontainers.NewStorageContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageContainerName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `StorageContainersClient.List`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `StorageContainersClient.ListAll`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListAll(ctx, id)` can be used to do batched pagination
items, err := client.ListAllComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `StorageContainersClient.Update`

```go
ctx := context.TODO()
id := storagecontainers.NewStorageContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageContainerName")

payload := storagecontainers.StorageContainersUpdateRequest{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


## `github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/syncgroupresource` Documentation

The `syncgroupresource` SDK allows for interaction with the Azure Resource Manager Service `storagesync` (API Version `2020-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/syncgroupresource"
```


### Client Initialization

```go
client := syncgroupresource.NewSyncGroupResourceClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SyncGroupResourceClient.SyncGroupsCreate`

```go
ctx := context.TODO()
id := syncgroupresource.NewSyncGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageSyncServiceValue", "syncGroupValue")

payload := syncgroupresource.SyncGroupCreateParameters{
	// ...
}


read, err := client.SyncGroupsCreate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SyncGroupResourceClient.SyncGroupsDelete`

```go
ctx := context.TODO()
id := syncgroupresource.NewSyncGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageSyncServiceValue", "syncGroupValue")

read, err := client.SyncGroupsDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SyncGroupResourceClient.SyncGroupsGet`

```go
ctx := context.TODO()
id := syncgroupresource.NewSyncGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageSyncServiceValue", "syncGroupValue")

read, err := client.SyncGroupsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SyncGroupResourceClient.SyncGroupsListByStorageSyncService`

```go
ctx := context.TODO()
id := syncgroupresource.NewStorageSyncServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageSyncServiceValue")

read, err := client.SyncGroupsListByStorageSyncService(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

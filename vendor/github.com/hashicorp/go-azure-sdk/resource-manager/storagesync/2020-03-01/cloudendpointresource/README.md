
## `github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/cloudendpointresource` Documentation

The `cloudendpointresource` SDK allows for interaction with Azure Resource Manager `storagesync` (API Version `2020-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/cloudendpointresource"
```


### Client Initialization

```go
client := cloudendpointresource.NewCloudEndpointResourceClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `CloudEndpointResourceClient.CloudEndpointsCreate`

```go
ctx := context.TODO()
id := cloudendpointresource.NewCloudEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageSyncServiceName", "syncGroupName", "cloudEndpointName")

payload := cloudendpointresource.CloudEndpointCreateParameters{
	// ...
}


if err := client.CloudEndpointsCreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CloudEndpointResourceClient.CloudEndpointsDelete`

```go
ctx := context.TODO()
id := cloudendpointresource.NewCloudEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageSyncServiceName", "syncGroupName", "cloudEndpointName")

if err := client.CloudEndpointsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CloudEndpointResourceClient.CloudEndpointsGet`

```go
ctx := context.TODO()
id := cloudendpointresource.NewCloudEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageSyncServiceName", "syncGroupName", "cloudEndpointName")

read, err := client.CloudEndpointsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CloudEndpointResourceClient.CloudEndpointsListBySyncGroup`

```go
ctx := context.TODO()
id := cloudendpointresource.NewSyncGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageSyncServiceName", "syncGroupName")

read, err := client.CloudEndpointsListBySyncGroup(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CloudEndpointResourceClient.CloudEndpointsPostBackup`

```go
ctx := context.TODO()
id := cloudendpointresource.NewCloudEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageSyncServiceName", "syncGroupName", "cloudEndpointName")

payload := cloudendpointresource.BackupRequest{
	// ...
}


if err := client.CloudEndpointsPostBackupThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CloudEndpointResourceClient.CloudEndpointsPostRestore`

```go
ctx := context.TODO()
id := cloudendpointresource.NewCloudEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageSyncServiceName", "syncGroupName", "cloudEndpointName")

payload := cloudendpointresource.PostRestoreRequest{
	// ...
}


if err := client.CloudEndpointsPostRestoreThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CloudEndpointResourceClient.CloudEndpointsPreBackup`

```go
ctx := context.TODO()
id := cloudendpointresource.NewCloudEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageSyncServiceName", "syncGroupName", "cloudEndpointName")

payload := cloudendpointresource.BackupRequest{
	// ...
}


if err := client.CloudEndpointsPreBackupThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CloudEndpointResourceClient.CloudEndpointsPreRestore`

```go
ctx := context.TODO()
id := cloudendpointresource.NewCloudEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageSyncServiceName", "syncGroupName", "cloudEndpointName")

payload := cloudendpointresource.PreRestoreRequest{
	// ...
}


if err := client.CloudEndpointsPreRestoreThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CloudEndpointResourceClient.CloudEndpointsTriggerChangeDetection`

```go
ctx := context.TODO()
id := cloudendpointresource.NewCloudEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageSyncServiceName", "syncGroupName", "cloudEndpointName")

payload := cloudendpointresource.TriggerChangeDetectionParameters{
	// ...
}


if err := client.CloudEndpointsTriggerChangeDetectionThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CloudEndpointResourceClient.CloudEndpointsrestoreheartbeat`

```go
ctx := context.TODO()
id := cloudendpointresource.NewCloudEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageSyncServiceName", "syncGroupName", "cloudEndpointName")

read, err := client.CloudEndpointsrestoreheartbeat(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


## `github.com/hashicorp/go-azure-sdk/resource-manager/hardwaresecuritymodules/2025-03-31/cloudhsmclusters` Documentation

The `cloudhsmclusters` SDK allows for interaction with Azure Resource Manager `hardwaresecuritymodules` (API Version `2025-03-31`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/hardwaresecuritymodules/2025-03-31/cloudhsmclusters"
```


### Client Initialization

```go
client := cloudhsmclusters.NewCloudHsmClustersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `CloudHsmClustersClient.Backup`

```go
ctx := context.TODO()
id := cloudhsmclusters.NewCloudHsmClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cloudHsmClusterName")

payload := cloudhsmclusters.BackupRestoreRequestBaseProperties{
	// ...
}


if err := client.BackupThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CloudHsmClustersClient.CloudHsmClusterBackupStatusGet`

```go
ctx := context.TODO()
id := cloudhsmclusters.NewBackupOperationStatusID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cloudHsmClusterName", "jobId")

read, err := client.CloudHsmClusterBackupStatusGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CloudHsmClustersClient.CloudHsmClusterPrivateLinkResourcesListByCloudHsmCluster`

```go
ctx := context.TODO()
id := cloudhsmclusters.NewCloudHsmClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cloudHsmClusterName")

// alternatively `client.CloudHsmClusterPrivateLinkResourcesListByCloudHsmCluster(ctx, id)` can be used to do batched pagination
items, err := client.CloudHsmClusterPrivateLinkResourcesListByCloudHsmClusterComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `CloudHsmClustersClient.CloudHsmClusterRestoreStatusGet`

```go
ctx := context.TODO()
id := cloudhsmclusters.NewRestoreOperationStatusID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cloudHsmClusterName", "jobId")

read, err := client.CloudHsmClusterRestoreStatusGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CloudHsmClustersClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := cloudhsmclusters.NewCloudHsmClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cloudHsmClusterName")

payload := cloudhsmclusters.CloudHsmCluster{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CloudHsmClustersClient.Delete`

```go
ctx := context.TODO()
id := cloudhsmclusters.NewCloudHsmClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cloudHsmClusterName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CloudHsmClustersClient.Get`

```go
ctx := context.TODO()
id := cloudhsmclusters.NewCloudHsmClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cloudHsmClusterName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CloudHsmClustersClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `CloudHsmClustersClient.ListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `CloudHsmClustersClient.Restore`

```go
ctx := context.TODO()
id := cloudhsmclusters.NewCloudHsmClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cloudHsmClusterName")

payload := cloudhsmclusters.RestoreRequestProperties{
	// ...
}


if err := client.RestoreThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CloudHsmClustersClient.Update`

```go
ctx := context.TODO()
id := cloudhsmclusters.NewCloudHsmClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cloudHsmClusterName")

payload := cloudhsmclusters.CloudHsmClusterPatchParameters{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CloudHsmClustersClient.ValidateBackupProperties`

```go
ctx := context.TODO()
id := cloudhsmclusters.NewCloudHsmClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cloudHsmClusterName")

payload := cloudhsmclusters.BackupRestoreRequestBaseProperties{
	// ...
}


if err := client.ValidateBackupPropertiesThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CloudHsmClustersClient.ValidateRestoreProperties`

```go
ctx := context.TODO()
id := cloudhsmclusters.NewCloudHsmClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cloudHsmClusterName")

payload := cloudhsmclusters.RestoreRequestProperties{
	// ...
}


if err := client.ValidateRestorePropertiesThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

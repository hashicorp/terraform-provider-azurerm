
## `github.com/hashicorp/go-azure-sdk/resource-manager/elasticsan/2023-01-01/snapshots` Documentation

The `snapshots` SDK allows for interaction with Azure Resource Manager `elasticsan` (API Version `2023-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/elasticsan/2023-01-01/snapshots"
```


### Client Initialization

```go
client := snapshots.NewSnapshotsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SnapshotsClient.VolumeSnapshotsCreate`

```go
ctx := context.TODO()
id := snapshots.NewSnapshotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "elasticSanName", "volumeGroupName", "snapshotName")

payload := snapshots.Snapshot{
	// ...
}


if err := client.VolumeSnapshotsCreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `SnapshotsClient.VolumeSnapshotsDelete`

```go
ctx := context.TODO()
id := snapshots.NewSnapshotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "elasticSanName", "volumeGroupName", "snapshotName")

if err := client.VolumeSnapshotsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `SnapshotsClient.VolumeSnapshotsGet`

```go
ctx := context.TODO()
id := snapshots.NewSnapshotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "elasticSanName", "volumeGroupName", "snapshotName")

read, err := client.VolumeSnapshotsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SnapshotsClient.VolumeSnapshotsListByVolumeGroup`

```go
ctx := context.TODO()
id := snapshots.NewVolumeGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "elasticSanName", "volumeGroupName")

// alternatively `client.VolumeSnapshotsListByVolumeGroup(ctx, id, snapshots.DefaultVolumeSnapshotsListByVolumeGroupOperationOptions())` can be used to do batched pagination
items, err := client.VolumeSnapshotsListByVolumeGroupComplete(ctx, id, snapshots.DefaultVolumeSnapshotsListByVolumeGroupOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

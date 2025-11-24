
## `github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-06-01/volumes` Documentation

The `volumes` SDK allows for interaction with Azure Resource Manager `netapp` (API Version `2025-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-06-01/volumes"
```


### Client Initialization

```go
client := volumes.NewVolumesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `VolumesClient.AuthorizeExternalReplication`

```go
ctx := context.TODO()
id := volumes.NewVolumeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "capacityPoolName", "volumeName")

if err := client.AuthorizeExternalReplicationThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VolumesClient.AuthorizeReplication`

```go
ctx := context.TODO()
id := volumes.NewVolumeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "capacityPoolName", "volumeName")

payload := volumes.AuthorizeRequest{
	// ...
}


if err := client.AuthorizeReplicationThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VolumesClient.BackupsGetLatestStatus`

```go
ctx := context.TODO()
id := volumes.NewVolumeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "capacityPoolName", "volumeName")

read, err := client.BackupsGetLatestStatus(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VolumesClient.BackupsGetVolumeLatestRestoreStatus`

```go
ctx := context.TODO()
id := volumes.NewVolumeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "capacityPoolName", "volumeName")

read, err := client.BackupsGetVolumeLatestRestoreStatus(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VolumesClient.BackupsUnderVolumeMigrateBackups`

```go
ctx := context.TODO()
id := volumes.NewVolumeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "capacityPoolName", "volumeName")

payload := volumes.BackupsMigrationRequest{
	// ...
}


if err := client.BackupsUnderVolumeMigrateBackupsThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VolumesClient.BreakFileLocks`

```go
ctx := context.TODO()
id := volumes.NewVolumeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "capacityPoolName", "volumeName")

payload := volumes.BreakFileLocksRequest{
	// ...
}


if err := client.BreakFileLocksThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VolumesClient.BreakReplication`

```go
ctx := context.TODO()
id := volumes.NewVolumeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "capacityPoolName", "volumeName")

payload := volumes.BreakReplicationRequest{
	// ...
}


if err := client.BreakReplicationThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VolumesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := volumes.NewVolumeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "capacityPoolName", "volumeName")

payload := volumes.Volume{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VolumesClient.Delete`

```go
ctx := context.TODO()
id := volumes.NewVolumeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "capacityPoolName", "volumeName")

if err := client.DeleteThenPoll(ctx, id, volumes.DefaultDeleteOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `VolumesClient.DeleteReplication`

```go
ctx := context.TODO()
id := volumes.NewVolumeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "capacityPoolName", "volumeName")

if err := client.DeleteReplicationThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VolumesClient.FinalizeExternalReplication`

```go
ctx := context.TODO()
id := volumes.NewVolumeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "capacityPoolName", "volumeName")

if err := client.FinalizeExternalReplicationThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VolumesClient.FinalizeRelocation`

```go
ctx := context.TODO()
id := volumes.NewVolumeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "capacityPoolName", "volumeName")

if err := client.FinalizeRelocationThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VolumesClient.Get`

```go
ctx := context.TODO()
id := volumes.NewVolumeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "capacityPoolName", "volumeName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VolumesClient.List`

```go
ctx := context.TODO()
id := volumes.NewCapacityPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "capacityPoolName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VolumesClient.ListGetGroupIdListForLdapUser`

```go
ctx := context.TODO()
id := volumes.NewVolumeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "capacityPoolName", "volumeName")

payload := volumes.GetGroupIdListForLdapUserRequest{
	// ...
}


if err := client.ListGetGroupIdListForLdapUserThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VolumesClient.ListReplications`

```go
ctx := context.TODO()
id := volumes.NewVolumeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "capacityPoolName", "volumeName")

// alternatively `client.ListReplications(ctx, id)` can be used to do batched pagination
items, err := client.ListReplicationsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VolumesClient.PeerExternalCluster`

```go
ctx := context.TODO()
id := volumes.NewVolumeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "capacityPoolName", "volumeName")

payload := volumes.PeerClusterForVolumeMigrationRequest{
	// ...
}


if err := client.PeerExternalClusterThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VolumesClient.PerformReplicationTransfer`

```go
ctx := context.TODO()
id := volumes.NewVolumeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "capacityPoolName", "volumeName")

if err := client.PerformReplicationTransferThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VolumesClient.PoolChange`

```go
ctx := context.TODO()
id := volumes.NewVolumeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "capacityPoolName", "volumeName")

payload := volumes.PoolChangeRequest{
	// ...
}


if err := client.PoolChangeThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VolumesClient.PopulateAvailabilityZone`

```go
ctx := context.TODO()
id := volumes.NewVolumeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "capacityPoolName", "volumeName")

if err := client.PopulateAvailabilityZoneThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VolumesClient.ReInitializeReplication`

```go
ctx := context.TODO()
id := volumes.NewVolumeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "capacityPoolName", "volumeName")

if err := client.ReInitializeReplicationThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VolumesClient.ReestablishReplication`

```go
ctx := context.TODO()
id := volumes.NewVolumeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "capacityPoolName", "volumeName")

payload := volumes.ReestablishReplicationRequest{
	// ...
}


if err := client.ReestablishReplicationThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VolumesClient.Relocate`

```go
ctx := context.TODO()
id := volumes.NewVolumeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "capacityPoolName", "volumeName")

payload := volumes.RelocateVolumeRequest{
	// ...
}


if err := client.RelocateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VolumesClient.ReplicationStatus`

```go
ctx := context.TODO()
id := volumes.NewVolumeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "capacityPoolName", "volumeName")

read, err := client.ReplicationStatus(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VolumesClient.ResetCifsPassword`

```go
ctx := context.TODO()
id := volumes.NewVolumeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "capacityPoolName", "volumeName")

if err := client.ResetCifsPasswordThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VolumesClient.ResyncReplication`

```go
ctx := context.TODO()
id := volumes.NewVolumeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "capacityPoolName", "volumeName")

if err := client.ResyncReplicationThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VolumesClient.Revert`

```go
ctx := context.TODO()
id := volumes.NewVolumeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "capacityPoolName", "volumeName")

payload := volumes.VolumeRevert{
	// ...
}


if err := client.RevertThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VolumesClient.RevertRelocation`

```go
ctx := context.TODO()
id := volumes.NewVolumeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "capacityPoolName", "volumeName")

if err := client.RevertRelocationThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VolumesClient.SplitCloneFromParent`

```go
ctx := context.TODO()
id := volumes.NewVolumeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "capacityPoolName", "volumeName")

if err := client.SplitCloneFromParentThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VolumesClient.Update`

```go
ctx := context.TODO()
id := volumes.NewVolumeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "capacityPoolName", "volumeName")

payload := volumes.VolumePatch{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

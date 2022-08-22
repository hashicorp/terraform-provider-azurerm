
## `github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2021-10-01/snapshots` Documentation

The `snapshots` SDK allows for interaction with the Azure Resource Manager Service `netapp` (API Version `2021-10-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2021-10-01/snapshots"
```


### Client Initialization

```go
client := snapshots.NewSnapshotsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SnapshotsClient.Create`

```go
ctx := context.TODO()
id := snapshots.NewSnapshotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue", "poolValue", "volumeValue", "snapshotValue")

payload := snapshots.Snapshot{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `SnapshotsClient.Delete`

```go
ctx := context.TODO()
id := snapshots.NewSnapshotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue", "poolValue", "volumeValue", "snapshotValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `SnapshotsClient.Get`

```go
ctx := context.TODO()
id := snapshots.NewSnapshotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue", "poolValue", "volumeValue", "snapshotValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SnapshotsClient.List`

```go
ctx := context.TODO()
id := snapshots.NewVolumeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue", "poolValue", "volumeValue")

read, err := client.List(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SnapshotsClient.RestoreFiles`

```go
ctx := context.TODO()
id := snapshots.NewSnapshotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue", "poolValue", "volumeValue", "snapshotValue")

payload := snapshots.SnapshotRestoreFiles{
	// ...
}


if err := client.RestoreFilesThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `SnapshotsClient.Update`

```go
ctx := context.TODO()
id := snapshots.NewSnapshotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue", "poolValue", "volumeValue", "snapshotValue")
var payload interface{}

if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

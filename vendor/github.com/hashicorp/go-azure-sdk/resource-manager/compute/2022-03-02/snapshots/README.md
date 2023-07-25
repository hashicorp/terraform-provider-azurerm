
## `github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/snapshots` Documentation

The `snapshots` SDK allows for interaction with the Azure Resource Manager Service `compute` (API Version `2022-03-02`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/snapshots"
```


### Client Initialization

```go
client := snapshots.NewSnapshotsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SnapshotsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := snapshots.NewSnapshotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "snapshotValue")

payload := snapshots.Snapshot{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `SnapshotsClient.Delete`

```go
ctx := context.TODO()
id := snapshots.NewSnapshotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "snapshotValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `SnapshotsClient.Get`

```go
ctx := context.TODO()
id := snapshots.NewSnapshotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "snapshotValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SnapshotsClient.GrantAccess`

```go
ctx := context.TODO()
id := snapshots.NewSnapshotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "snapshotValue")

payload := snapshots.GrantAccessData{
	// ...
}


if err := client.GrantAccessThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `SnapshotsClient.List`

```go
ctx := context.TODO()
id := snapshots.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `SnapshotsClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := snapshots.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `SnapshotsClient.RevokeAccess`

```go
ctx := context.TODO()
id := snapshots.NewSnapshotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "snapshotValue")

if err := client.RevokeAccessThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `SnapshotsClient.Update`

```go
ctx := context.TODO()
id := snapshots.NewSnapshotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "snapshotValue")

payload := snapshots.SnapshotUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

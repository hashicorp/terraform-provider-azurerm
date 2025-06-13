
## `github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-03-02-preview/managedclustersnapshots` Documentation

The `managedclustersnapshots` SDK allows for interaction with Azure Resource Manager `containerservice` (API Version `2023-03-02-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-03-02-preview/managedclustersnapshots"
```


### Client Initialization

```go
client := managedclustersnapshots.NewManagedClusterSnapshotsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ManagedClusterSnapshotsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := managedclustersnapshots.NewManagedClusterSnapshotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedClusterSnapshotName")

payload := managedclustersnapshots.ManagedClusterSnapshot{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedClusterSnapshotsClient.Delete`

```go
ctx := context.TODO()
id := managedclustersnapshots.NewManagedClusterSnapshotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedClusterSnapshotName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedClusterSnapshotsClient.Get`

```go
ctx := context.TODO()
id := managedclustersnapshots.NewManagedClusterSnapshotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedClusterSnapshotName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedClusterSnapshotsClient.List`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ManagedClusterSnapshotsClient.ListByResourceGroup`

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


### Example Usage: `ManagedClusterSnapshotsClient.UpdateTags`

```go
ctx := context.TODO()
id := managedclustersnapshots.NewManagedClusterSnapshotID("12345678-1234-9876-4563-123456789012", "example-resource-group", "managedClusterSnapshotName")

payload := managedclustersnapshots.TagsObject{
	// ...
}


read, err := client.UpdateTags(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

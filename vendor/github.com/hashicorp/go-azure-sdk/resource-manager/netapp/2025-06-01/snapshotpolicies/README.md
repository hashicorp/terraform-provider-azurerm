
## `github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-06-01/snapshotpolicies` Documentation

The `snapshotpolicies` SDK allows for interaction with Azure Resource Manager `netapp` (API Version `2025-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-06-01/snapshotpolicies"
```


### Client Initialization

```go
client := snapshotpolicies.NewSnapshotPoliciesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SnapshotPoliciesClient.Create`

```go
ctx := context.TODO()
id := snapshotpolicies.NewSnapshotPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "snapshotPolicyName")

payload := snapshotpolicies.SnapshotPolicy{
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


### Example Usage: `SnapshotPoliciesClient.Delete`

```go
ctx := context.TODO()
id := snapshotpolicies.NewSnapshotPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "snapshotPolicyName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `SnapshotPoliciesClient.Get`

```go
ctx := context.TODO()
id := snapshotpolicies.NewSnapshotPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "snapshotPolicyName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SnapshotPoliciesClient.List`

```go
ctx := context.TODO()
id := snapshotpolicies.NewNetAppAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `SnapshotPoliciesClient.ListVolumes`

```go
ctx := context.TODO()
id := snapshotpolicies.NewSnapshotPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "snapshotPolicyName")

// alternatively `client.ListVolumes(ctx, id)` can be used to do batched pagination
items, err := client.ListVolumesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `SnapshotPoliciesClient.Update`

```go
ctx := context.TODO()
id := snapshotpolicies.NewSnapshotPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "snapshotPolicyName")

payload := snapshotpolicies.SnapshotPolicyPatch{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

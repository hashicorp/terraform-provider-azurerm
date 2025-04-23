
## `github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-01-01/snapshotpolicy` Documentation

The `snapshotpolicy` SDK allows for interaction with Azure Resource Manager `netapp` (API Version `2025-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-01-01/snapshotpolicy"
```


### Client Initialization

```go
client := snapshotpolicy.NewSnapshotPolicyClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SnapshotPolicyClient.SnapshotPoliciesCreate`

```go
ctx := context.TODO()
id := snapshotpolicy.NewSnapshotPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "snapshotPolicyName")

payload := snapshotpolicy.SnapshotPolicy{
	// ...
}


read, err := client.SnapshotPoliciesCreate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SnapshotPolicyClient.SnapshotPoliciesDelete`

```go
ctx := context.TODO()
id := snapshotpolicy.NewSnapshotPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "snapshotPolicyName")

if err := client.SnapshotPoliciesDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `SnapshotPolicyClient.SnapshotPoliciesGet`

```go
ctx := context.TODO()
id := snapshotpolicy.NewSnapshotPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "snapshotPolicyName")

read, err := client.SnapshotPoliciesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SnapshotPolicyClient.SnapshotPoliciesList`

```go
ctx := context.TODO()
id := snapshotpolicy.NewNetAppAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName")

read, err := client.SnapshotPoliciesList(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SnapshotPolicyClient.SnapshotPoliciesUpdate`

```go
ctx := context.TODO()
id := snapshotpolicy.NewSnapshotPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "snapshotPolicyName")

payload := snapshotpolicy.SnapshotPolicyPatch{
	// ...
}


if err := client.SnapshotPoliciesUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

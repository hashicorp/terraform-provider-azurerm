
## `github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-06-01/backuppolicies` Documentation

The `backuppolicies` SDK allows for interaction with Azure Resource Manager `netapp` (API Version `2025-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-06-01/backuppolicies"
```


### Client Initialization

```go
client := backuppolicies.NewBackupPoliciesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `BackupPoliciesClient.Create`

```go
ctx := context.TODO()
id := backuppolicies.NewBackupPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "backupPolicyName")

payload := backuppolicies.BackupPolicy{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `BackupPoliciesClient.Delete`

```go
ctx := context.TODO()
id := backuppolicies.NewBackupPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "backupPolicyName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `BackupPoliciesClient.Get`

```go
ctx := context.TODO()
id := backuppolicies.NewBackupPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "backupPolicyName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BackupPoliciesClient.List`

```go
ctx := context.TODO()
id := backuppolicies.NewNetAppAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `BackupPoliciesClient.Update`

```go
ctx := context.TODO()
id := backuppolicies.NewBackupPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "backupPolicyName")

payload := backuppolicies.BackupPolicyPatch{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

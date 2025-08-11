
## `github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-01-01/backuppolicy` Documentation

The `backuppolicy` SDK allows for interaction with Azure Resource Manager `netapp` (API Version `2025-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-01-01/backuppolicy"
```


### Client Initialization

```go
client := backuppolicy.NewBackupPolicyClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `BackupPolicyClient.BackupPoliciesCreate`

```go
ctx := context.TODO()
id := backuppolicy.NewBackupPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "backupPolicyName")

payload := backuppolicy.BackupPolicy{
	// ...
}


if err := client.BackupPoliciesCreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `BackupPolicyClient.BackupPoliciesDelete`

```go
ctx := context.TODO()
id := backuppolicy.NewBackupPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "backupPolicyName")

if err := client.BackupPoliciesDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `BackupPolicyClient.BackupPoliciesGet`

```go
ctx := context.TODO()
id := backuppolicy.NewBackupPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "backupPolicyName")

read, err := client.BackupPoliciesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BackupPolicyClient.BackupPoliciesList`

```go
ctx := context.TODO()
id := backuppolicy.NewNetAppAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName")

read, err := client.BackupPoliciesList(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BackupPolicyClient.BackupPoliciesUpdate`

```go
ctx := context.TODO()
id := backuppolicy.NewBackupPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "backupPolicyName")

payload := backuppolicy.BackupPolicyPatch{
	// ...
}


if err := client.BackupPoliciesUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

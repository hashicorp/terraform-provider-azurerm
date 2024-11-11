
## `github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2024-04-01/backuppolicies` Documentation

The `backuppolicies` SDK allows for interaction with Azure Resource Manager `dataprotection` (API Version `2024-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2024-04-01/backuppolicies"
```


### Client Initialization

```go
client := backuppolicies.NewBackupPoliciesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `BackupPoliciesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := backuppolicies.NewBackupPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultName", "backupPolicyName")

payload := backuppolicies.BaseBackupPolicyResource{
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


### Example Usage: `BackupPoliciesClient.Delete`

```go
ctx := context.TODO()
id := backuppolicies.NewBackupPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultName", "backupPolicyName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BackupPoliciesClient.Get`

```go
ctx := context.TODO()
id := backuppolicies.NewBackupPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultName", "backupPolicyName")

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
id := backuppolicies.NewBackupVaultID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

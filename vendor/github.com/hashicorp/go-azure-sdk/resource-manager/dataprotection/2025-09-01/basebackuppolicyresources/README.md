
## `github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-09-01/basebackuppolicyresources` Documentation

The `basebackuppolicyresources` SDK allows for interaction with Azure Resource Manager `dataprotection` (API Version `2025-09-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-09-01/basebackuppolicyresources"
```


### Client Initialization

```go
client := basebackuppolicyresources.NewBaseBackupPolicyResourcesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `BaseBackupPolicyResourcesClient.BackupPoliciesCreateOrUpdate`

```go
ctx := context.TODO()
id := basebackuppolicyresources.NewBackupPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultName", "backupPolicyName")

payload := basebackuppolicyresources.BaseBackupPolicyResource{
	// ...
}


read, err := client.BackupPoliciesCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BaseBackupPolicyResourcesClient.BackupPoliciesDelete`

```go
ctx := context.TODO()
id := basebackuppolicyresources.NewBackupPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultName", "backupPolicyName")

read, err := client.BackupPoliciesDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BaseBackupPolicyResourcesClient.BackupPoliciesGet`

```go
ctx := context.TODO()
id := basebackuppolicyresources.NewBackupPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultName", "backupPolicyName")

read, err := client.BackupPoliciesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BaseBackupPolicyResourcesClient.BackupPoliciesList`

```go
ctx := context.TODO()
id := basebackuppolicyresources.NewBackupVaultID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultName")

// alternatively `client.BackupPoliciesList(ctx, id)` can be used to do batched pagination
items, err := client.BackupPoliciesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

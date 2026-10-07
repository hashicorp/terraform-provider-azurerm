
## `github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/backupvaultresources` Documentation

The `backupvaultresources` SDK allows for interaction with Azure Resource Manager `dataprotection` (API Version `2025-07-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/backupvaultresources"
```


### Client Initialization

```go
client := backupvaultresources.NewBackupVaultResourcesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `BackupVaultResourcesClient.BackupInstancesValidateForBackup`

```go
ctx := context.TODO()
id := backupvaultresources.NewBackupVaultID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultName")

payload := backupvaultresources.ValidateForBackupRequest{
	// ...
}


if err := client.BackupInstancesValidateForBackupThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `BackupVaultResourcesClient.BackupVaultsCreateOrUpdate`

```go
ctx := context.TODO()
id := backupvaultresources.NewBackupVaultID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultName")

payload := backupvaultresources.BackupVaultResource{
	// ...
}


if err := client.BackupVaultsCreateOrUpdateThenPoll(ctx, id, payload, backupvaultresources.DefaultBackupVaultsCreateOrUpdateOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `BackupVaultResourcesClient.BackupVaultsDelete`

```go
ctx := context.TODO()
id := backupvaultresources.NewBackupVaultID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultName")

if err := client.BackupVaultsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `BackupVaultResourcesClient.BackupVaultsGet`

```go
ctx := context.TODO()
id := backupvaultresources.NewBackupVaultID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultName")

read, err := client.BackupVaultsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BackupVaultResourcesClient.BackupVaultsGetInSubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.BackupVaultsGetInSubscription(ctx, id)` can be used to do batched pagination
items, err := client.BackupVaultsGetInSubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `BackupVaultResourcesClient.BackupVaultsUpdate`

```go
ctx := context.TODO()
id := backupvaultresources.NewBackupVaultID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultName")

payload := backupvaultresources.PatchResourceRequestInput{
	// ...
}


if err := client.BackupVaultsUpdateThenPoll(ctx, id, payload, backupvaultresources.DefaultBackupVaultsUpdateOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `BackupVaultResourcesClient.ExportJobsOperationResultGet`

```go
ctx := context.TODO()
id := backupvaultresources.NewOperationIdID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultName", "operationId")

read, err := client.ExportJobsOperationResultGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BackupVaultResourcesClient.ExportJobsTrigger`

```go
ctx := context.TODO()
id := backupvaultresources.NewBackupVaultID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultName")

if err := client.ExportJobsTriggerThenPoll(ctx, id); err != nil {
	// handle the error
}
```

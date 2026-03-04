
## `github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-09-01/backupinstanceresources` Documentation

The `backupinstanceresources` SDK allows for interaction with Azure Resource Manager `dataprotection` (API Version `2025-09-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-09-01/backupinstanceresources"
```


### Client Initialization

```go
client := backupinstanceresources.NewBackupInstanceResourcesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `BackupInstanceResourcesClient.BackupInstancesAdhocBackup`

```go
ctx := context.TODO()
id := backupinstanceresources.NewBackupInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultName", "backupInstanceName")

payload := backupinstanceresources.TriggerBackupRequest{
	// ...
}


if err := client.BackupInstancesAdhocBackupThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `BackupInstanceResourcesClient.BackupInstancesCreateOrUpdate`

```go
ctx := context.TODO()
id := backupinstanceresources.NewBackupInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultName", "backupInstanceName")

payload := backupinstanceresources.BackupInstanceResource{
	// ...
}


if err := client.BackupInstancesCreateOrUpdateThenPoll(ctx, id, payload, backupinstanceresources.DefaultBackupInstancesCreateOrUpdateOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `BackupInstanceResourcesClient.BackupInstancesDelete`

```go
ctx := context.TODO()
id := backupinstanceresources.NewBackupInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultName", "backupInstanceName")

if err := client.BackupInstancesDeleteThenPoll(ctx, id, backupinstanceresources.DefaultBackupInstancesDeleteOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `BackupInstanceResourcesClient.BackupInstancesGet`

```go
ctx := context.TODO()
id := backupinstanceresources.NewBackupInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultName", "backupInstanceName")

read, err := client.BackupInstancesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BackupInstanceResourcesClient.BackupInstancesResumeBackups`

```go
ctx := context.TODO()
id := backupinstanceresources.NewBackupInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultName", "backupInstanceName")

if err := client.BackupInstancesResumeBackupsThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `BackupInstanceResourcesClient.BackupInstancesResumeProtection`

```go
ctx := context.TODO()
id := backupinstanceresources.NewBackupInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultName", "backupInstanceName")

if err := client.BackupInstancesResumeProtectionThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `BackupInstanceResourcesClient.BackupInstancesStopProtection`

```go
ctx := context.TODO()
id := backupinstanceresources.NewBackupInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultName", "backupInstanceName")

payload := backupinstanceresources.StopProtectionRequest{
	// ...
}


if err := client.BackupInstancesStopProtectionThenPoll(ctx, id, payload, backupinstanceresources.DefaultBackupInstancesStopProtectionOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `BackupInstanceResourcesClient.BackupInstancesSuspendBackups`

```go
ctx := context.TODO()
id := backupinstanceresources.NewBackupInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultName", "backupInstanceName")

payload := backupinstanceresources.SuspendBackupRequest{
	// ...
}


if err := client.BackupInstancesSuspendBackupsThenPoll(ctx, id, payload, backupinstanceresources.DefaultBackupInstancesSuspendBackupsOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `BackupInstanceResourcesClient.BackupInstancesSyncBackupInstance`

```go
ctx := context.TODO()
id := backupinstanceresources.NewBackupInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultName", "backupInstanceName")

payload := backupinstanceresources.SyncBackupInstanceRequest{
	// ...
}


if err := client.BackupInstancesSyncBackupInstanceThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `BackupInstanceResourcesClient.BackupInstancesTriggerRehydrate`

```go
ctx := context.TODO()
id := backupinstanceresources.NewBackupInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultName", "backupInstanceName")

payload := backupinstanceresources.AzureBackupRehydrationRequest{
	// ...
}


if err := client.BackupInstancesTriggerRehydrateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `BackupInstanceResourcesClient.BackupInstancesTriggerRestore`

```go
ctx := context.TODO()
id := backupinstanceresources.NewBackupInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultName", "backupInstanceName")

payload := backupinstanceresources.AzureBackupRestoreRequest{
	// ...
}


if err := client.BackupInstancesTriggerRestoreThenPoll(ctx, id, payload, backupinstanceresources.DefaultBackupInstancesTriggerRestoreOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `BackupInstanceResourcesClient.BackupInstancesValidateForModifyBackup`

```go
ctx := context.TODO()
id := backupinstanceresources.NewBackupInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultName", "backupInstanceName")

payload := backupinstanceresources.ValidateForModifyBackupRequest{
	// ...
}


if err := client.BackupInstancesValidateForModifyBackupThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `BackupInstanceResourcesClient.BackupInstancesValidateForRestore`

```go
ctx := context.TODO()
id := backupinstanceresources.NewBackupInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultName", "backupInstanceName")

payload := backupinstanceresources.ValidateRestoreRequestObject{
	// ...
}


if err := client.BackupInstancesValidateForRestoreThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `BackupInstanceResourcesClient.RestorableTimeRangesFind`

```go
ctx := context.TODO()
id := backupinstanceresources.NewBackupInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultName", "backupInstanceName")

payload := backupinstanceresources.AzureBackupFindRestorableTimeRangesRequest{
	// ...
}


read, err := client.RestorableTimeRangesFind(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

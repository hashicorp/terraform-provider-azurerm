
## `github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2024-04-01/backupinstances` Documentation

The `backupinstances` SDK allows for interaction with Azure Resource Manager `dataprotection` (API Version `2024-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2024-04-01/backupinstances"
```


### Client Initialization

```go
client := backupinstances.NewBackupInstancesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `BackupInstancesClient.AdhocBackup`

```go
ctx := context.TODO()
id := backupinstances.NewBackupInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultName", "backupInstanceName")

payload := backupinstances.TriggerBackupRequest{
	// ...
}


if err := client.AdhocBackupThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `BackupInstancesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := backupinstances.NewBackupInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultName", "backupInstanceName")

payload := backupinstances.BackupInstanceResource{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload, backupinstances.DefaultCreateOrUpdateOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `BackupInstancesClient.Delete`

```go
ctx := context.TODO()
id := backupinstances.NewBackupInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultName", "backupInstanceName")

if err := client.DeleteThenPoll(ctx, id, backupinstances.DefaultDeleteOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `BackupInstancesClient.Get`

```go
ctx := context.TODO()
id := backupinstances.NewBackupInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultName", "backupInstanceName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BackupInstancesClient.List`

```go
ctx := context.TODO()
id := backupinstances.NewBackupVaultID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `BackupInstancesClient.ResumeBackups`

```go
ctx := context.TODO()
id := backupinstances.NewBackupInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultName", "backupInstanceName")

if err := client.ResumeBackupsThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `BackupInstancesClient.ResumeProtection`

```go
ctx := context.TODO()
id := backupinstances.NewBackupInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultName", "backupInstanceName")

if err := client.ResumeProtectionThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `BackupInstancesClient.StopProtection`

```go
ctx := context.TODO()
id := backupinstances.NewBackupInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultName", "backupInstanceName")

payload := backupinstances.StopProtectionRequest{
	// ...
}


if err := client.StopProtectionThenPoll(ctx, id, payload, backupinstances.DefaultStopProtectionOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `BackupInstancesClient.SuspendBackups`

```go
ctx := context.TODO()
id := backupinstances.NewBackupInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultName", "backupInstanceName")

payload := backupinstances.SuspendBackupRequest{
	// ...
}


if err := client.SuspendBackupsThenPoll(ctx, id, payload, backupinstances.DefaultSuspendBackupsOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `BackupInstancesClient.SyncBackupInstance`

```go
ctx := context.TODO()
id := backupinstances.NewBackupInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultName", "backupInstanceName")

payload := backupinstances.SyncBackupInstanceRequest{
	// ...
}


if err := client.SyncBackupInstanceThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `BackupInstancesClient.TriggerCrossRegionRestore`

```go
ctx := context.TODO()
id := backupinstances.NewProviderLocationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "locationName")

payload := backupinstances.CrossRegionRestoreRequestObject{
	// ...
}


if err := client.TriggerCrossRegionRestoreThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `BackupInstancesClient.TriggerRehydrate`

```go
ctx := context.TODO()
id := backupinstances.NewBackupInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultName", "backupInstanceName")

payload := backupinstances.AzureBackupRehydrationRequest{
	// ...
}


if err := client.TriggerRehydrateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `BackupInstancesClient.TriggerRestore`

```go
ctx := context.TODO()
id := backupinstances.NewBackupInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultName", "backupInstanceName")

payload := backupinstances.AzureBackupRestoreRequest{
	// ...
}


if err := client.TriggerRestoreThenPoll(ctx, id, payload, backupinstances.DefaultTriggerRestoreOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `BackupInstancesClient.ValidateCrossRegionRestore`

```go
ctx := context.TODO()
id := backupinstances.NewProviderLocationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "locationName")

payload := backupinstances.ValidateCrossRegionRestoreRequestObject{
	// ...
}


if err := client.ValidateCrossRegionRestoreThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `BackupInstancesClient.ValidateForBackup`

```go
ctx := context.TODO()
id := backupinstances.NewBackupVaultID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultName")

payload := backupinstances.ValidateForBackupRequest{
	// ...
}


if err := client.ValidateForBackupThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `BackupInstancesClient.ValidateForRestore`

```go
ctx := context.TODO()
id := backupinstances.NewBackupInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultName", "backupInstanceName")

payload := backupinstances.ValidateRestoreRequestObject{
	// ...
}


if err := client.ValidateForRestoreThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

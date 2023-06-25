
## `github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2022-04-01/backupinstances` Documentation

The `backupinstances` SDK allows for interaction with the Azure Resource Manager Service `dataprotection` (API Version `2022-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2022-04-01/backupinstances"
```


### Client Initialization

```go
client := backupinstances.NewBackupInstancesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `BackupInstancesClient.AdhocBackup`

```go
ctx := context.TODO()
id := backupinstances.NewBackupInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultValue", "backupInstanceValue")

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
id := backupinstances.NewBackupInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultValue", "backupInstanceValue")

payload := backupinstances.BackupInstanceResource{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `BackupInstancesClient.Delete`

```go
ctx := context.TODO()
id := backupinstances.NewBackupInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultValue", "backupInstanceValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `BackupInstancesClient.Get`

```go
ctx := context.TODO()
id := backupinstances.NewBackupInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultValue", "backupInstanceValue")

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
id := backupinstances.NewBackupVaultID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultValue")

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
id := backupinstances.NewBackupInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultValue", "backupInstanceValue")

if err := client.ResumeBackupsThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `BackupInstancesClient.ResumeProtection`

```go
ctx := context.TODO()
id := backupinstances.NewBackupInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultValue", "backupInstanceValue")

if err := client.ResumeProtectionThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `BackupInstancesClient.StopProtection`

```go
ctx := context.TODO()
id := backupinstances.NewBackupInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultValue", "backupInstanceValue")

if err := client.StopProtectionThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `BackupInstancesClient.SuspendBackups`

```go
ctx := context.TODO()
id := backupinstances.NewBackupInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultValue", "backupInstanceValue")

if err := client.SuspendBackupsThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `BackupInstancesClient.SyncBackupInstance`

```go
ctx := context.TODO()
id := backupinstances.NewBackupInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultValue", "backupInstanceValue")

payload := backupinstances.SyncBackupInstanceRequest{
	// ...
}


if err := client.SyncBackupInstanceThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `BackupInstancesClient.TriggerRehydrate`

```go
ctx := context.TODO()
id := backupinstances.NewBackupInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultValue", "backupInstanceValue")

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
id := backupinstances.NewBackupInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultValue", "backupInstanceValue")

payload := backupinstances.AzureBackupRestoreRequest{
	// ...
}


if err := client.TriggerRestoreThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `BackupInstancesClient.ValidateForBackup`

```go
ctx := context.TODO()
id := backupinstances.NewBackupVaultID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultValue")

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
id := backupinstances.NewBackupInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "backupVaultValue", "backupInstanceValue")

payload := backupinstances.ValidateRestoreRequestObject{
	// ...
}


if err := client.ValidateForRestoreThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

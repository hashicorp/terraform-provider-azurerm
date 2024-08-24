
## `github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2024-03-01/backups` Documentation

The `backups` SDK allows for interaction with the Azure Resource Manager Service `netapp` (API Version `2024-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2024-03-01/backups"
```


### Client Initialization

```go
client := backups.NewBackupsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `BackupsClient.Create`

```go
ctx := context.TODO()
id := backups.NewBackupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountValue", "backupVaultValue", "backupValue")

payload := backups.Backup{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `BackupsClient.Delete`

```go
ctx := context.TODO()
id := backups.NewBackupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountValue", "backupVaultValue", "backupValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `BackupsClient.Get`

```go
ctx := context.TODO()
id := backups.NewBackupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountValue", "backupVaultValue", "backupValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BackupsClient.GetLatestStatus`

```go
ctx := context.TODO()
id := backups.NewVolumeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountValue", "capacityPoolValue", "volumeValue")

read, err := client.GetLatestStatus(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BackupsClient.ListByVault`

```go
ctx := context.TODO()
id := backups.NewBackupVaultID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountValue", "backupVaultValue")

// alternatively `client.ListByVault(ctx, id, backups.DefaultListByVaultOperationOptions())` can be used to do batched pagination
items, err := client.ListByVaultComplete(ctx, id, backups.DefaultListByVaultOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `BackupsClient.UnderAccountMigrateBackups`

```go
ctx := context.TODO()
id := backups.NewNetAppAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountValue")

payload := backups.BackupsMigrationRequest{
	// ...
}


if err := client.UnderAccountMigrateBackupsThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `BackupsClient.UnderBackupVaultRestoreFiles`

```go
ctx := context.TODO()
id := backups.NewBackupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountValue", "backupVaultValue", "backupValue")

payload := backups.BackupRestoreFiles{
	// ...
}


if err := client.UnderBackupVaultRestoreFilesThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `BackupsClient.UnderVolumeMigrateBackups`

```go
ctx := context.TODO()
id := backups.NewVolumeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountValue", "capacityPoolValue", "volumeValue")

payload := backups.BackupsMigrationRequest{
	// ...
}


if err := client.UnderVolumeMigrateBackupsThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `BackupsClient.Update`

```go
ctx := context.TODO()
id := backups.NewBackupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountValue", "backupVaultValue", "backupValue")

payload := backups.BackupPatch{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

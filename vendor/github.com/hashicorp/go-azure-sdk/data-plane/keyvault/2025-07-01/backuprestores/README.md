
## `github.com/hashicorp/go-azure-sdk/data-plane/keyvault/2025-07-01/backuprestores` Documentation

The `backuprestores` SDK allows for interaction with <unknown source data type> `keyvault` (API Version `2025-07-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/data-plane/keyvault/2025-07-01/backuprestores"
```


### Client Initialization

```go
client := backuprestores.NewBackuprestoresClientWithBaseURI("")
client.Client.Authorizer = authorizer
```


### Example Usage: `BackuprestoresClient.FullBackup`

```go
ctx := context.TODO()

payload := backuprestores.SASTokenParameter{
	// ...
}


if err := client.FullBackupThenPoll(ctx, payload); err != nil {
	// handle the error
}
```


### Example Usage: `BackuprestoresClient.FullBackupStatus`

```go
ctx := context.TODO()
id := backuprestores.NewBackupID("jobId")

read, err := client.FullBackupStatus(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BackuprestoresClient.FullRestoreOperation`

```go
ctx := context.TODO()

payload := backuprestores.RestoreOperationParameters{
	// ...
}


if err := client.FullRestoreOperationThenPoll(ctx, payload); err != nil {
	// handle the error
}
```


### Example Usage: `BackuprestoresClient.PreFullBackup`

```go
ctx := context.TODO()

payload := backuprestores.PreBackupOperationParameters{
	// ...
}


if err := client.PreFullBackupThenPoll(ctx, payload); err != nil {
	// handle the error
}
```


### Example Usage: `BackuprestoresClient.PreFullRestoreOperation`

```go
ctx := context.TODO()

payload := backuprestores.PreRestoreOperationParameters{
	// ...
}


if err := client.PreFullRestoreOperationThenPoll(ctx, payload); err != nil {
	// handle the error
}
```


### Example Usage: `BackuprestoresClient.RestoreStatus`

```go
ctx := context.TODO()
id := backuprestores.NewRestoreID("jobId")

read, err := client.RestoreStatus(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BackuprestoresClient.SelectiveKeyRestoreOperation`

```go
ctx := context.TODO()
id := backuprestores.NewKeyID("keyName")

payload := backuprestores.SelectiveKeyRestoreOperationParameters{
	// ...
}


if err := client.SelectiveKeyRestoreOperationThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

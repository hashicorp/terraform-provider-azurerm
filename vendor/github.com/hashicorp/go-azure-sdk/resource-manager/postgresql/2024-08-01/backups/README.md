
## `github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2024-08-01/backups` Documentation

The `backups` SDK allows for interaction with Azure Resource Manager `postgresql` (API Version `2024-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2024-08-01/backups"
```


### Client Initialization

```go
client := backups.NewBackupsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `BackupsClient.Create`

```go
ctx := context.TODO()
id := backups.NewBackupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerName", "backupName")

if err := client.CreateThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `BackupsClient.Delete`

```go
ctx := context.TODO()
id := backups.NewBackupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerName", "backupName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `BackupsClient.Get`

```go
ctx := context.TODO()
id := backups.NewBackupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerName", "backupName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BackupsClient.ListByServer`

```go
ctx := context.TODO()
id := backups.NewFlexibleServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerName")

// alternatively `client.ListByServer(ctx, id)` can be used to do batched pagination
items, err := client.ListByServerComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

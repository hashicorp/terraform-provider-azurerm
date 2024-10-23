
## `github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/autonomousdatabasebackups` Documentation

The `autonomousdatabasebackups` SDK allows for interaction with Azure Resource Manager `oracledatabase` (API Version `2024-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/autonomousdatabasebackups"
```


### Client Initialization

```go
client := autonomousdatabasebackups.NewAutonomousDatabaseBackupsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AutonomousDatabaseBackupsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := autonomousdatabasebackups.NewAutonomousDatabaseBackupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "autonomousDatabaseName", "autonomousDatabaseBackupName")

payload := autonomousdatabasebackups.AutonomousDatabaseBackup{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AutonomousDatabaseBackupsClient.Delete`

```go
ctx := context.TODO()
id := autonomousdatabasebackups.NewAutonomousDatabaseBackupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "autonomousDatabaseName", "autonomousDatabaseBackupName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AutonomousDatabaseBackupsClient.Get`

```go
ctx := context.TODO()
id := autonomousdatabasebackups.NewAutonomousDatabaseBackupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "autonomousDatabaseName", "autonomousDatabaseBackupName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AutonomousDatabaseBackupsClient.ListByAutonomousDatabase`

```go
ctx := context.TODO()
id := autonomousdatabasebackups.NewAutonomousDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "autonomousDatabaseName")

// alternatively `client.ListByAutonomousDatabase(ctx, id)` can be used to do batched pagination
items, err := client.ListByAutonomousDatabaseComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AutonomousDatabaseBackupsClient.Update`

```go
ctx := context.TODO()
id := autonomousdatabasebackups.NewAutonomousDatabaseBackupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "autonomousDatabaseName", "autonomousDatabaseBackupName")

payload := autonomousdatabasebackups.AutonomousDatabaseBackupUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

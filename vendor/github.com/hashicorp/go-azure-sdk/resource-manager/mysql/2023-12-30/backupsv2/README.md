
## `github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/backupsv2` Documentation

The `backupsv2` SDK allows for interaction with Azure Resource Manager `mysql` (API Version `2023-12-30`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/backupsv2"
```


### Client Initialization

```go
client := backupsv2.NewBackupsV2ClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `BackupsV2Client.LongRunningBackupCreate`

```go
ctx := context.TODO()
id := backupsv2.NewBackupsV2ID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerName", "backupsV2Name")

payload := backupsv2.ServerBackupV2{
	// ...
}


if err := client.LongRunningBackupCreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `BackupsV2Client.LongRunningBackupsGet`

```go
ctx := context.TODO()
id := backupsv2.NewBackupsV2ID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerName", "backupsV2Name")

read, err := client.LongRunningBackupsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BackupsV2Client.LongRunningBackupsList`

```go
ctx := context.TODO()
id := backupsv2.NewFlexibleServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerName")

// alternatively `client.LongRunningBackupsList(ctx, id)` can be used to do batched pagination
items, err := client.LongRunningBackupsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

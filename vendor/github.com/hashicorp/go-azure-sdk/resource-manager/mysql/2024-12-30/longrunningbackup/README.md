
## `github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2024-12-30/longrunningbackup` Documentation

The `longrunningbackup` SDK allows for interaction with Azure Resource Manager `mysql` (API Version `2024-12-30`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2024-12-30/longrunningbackup"
```


### Client Initialization

```go
client := longrunningbackup.NewLongRunningBackupClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `LongRunningBackupClient.Create`

```go
ctx := context.TODO()
id := longrunningbackup.NewBackupsV2ID("12345678-1234-9876-4563-123456789012", "example-resource-group", "flexibleServerName", "backupsV2Name")

payload := longrunningbackup.ServerBackupV2{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

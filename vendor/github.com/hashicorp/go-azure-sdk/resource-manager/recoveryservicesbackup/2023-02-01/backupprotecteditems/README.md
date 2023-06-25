
## `github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/backupprotecteditems` Documentation

The `backupprotecteditems` SDK allows for interaction with the Azure Resource Manager Service `recoveryservicesbackup` (API Version `2023-02-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/backupprotecteditems"
```


### Client Initialization

```go
client := backupprotecteditems.NewBackupProtectedItemsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `BackupProtectedItemsClient.List`

```go
ctx := context.TODO()
id := backupprotecteditems.NewVaultID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue")

// alternatively `client.List(ctx, id, backupprotecteditems.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, backupprotecteditems.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

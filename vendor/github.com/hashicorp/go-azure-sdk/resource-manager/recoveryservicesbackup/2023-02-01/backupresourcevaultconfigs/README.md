
## `github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/backupresourcevaultconfigs` Documentation

The `backupresourcevaultconfigs` SDK allows for interaction with the Azure Resource Manager Service `recoveryservicesbackup` (API Version `2023-02-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/backupresourcevaultconfigs"
```


### Client Initialization

```go
client := backupresourcevaultconfigs.NewBackupResourceVaultConfigsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `BackupResourceVaultConfigsClient.Get`

```go
ctx := context.TODO()
id := backupresourcevaultconfigs.NewVaultID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BackupResourceVaultConfigsClient.Put`

```go
ctx := context.TODO()
id := backupresourcevaultconfigs.NewVaultID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue")

payload := backupresourcevaultconfigs.BackupResourceVaultConfigResource{
	// ...
}


read, err := client.Put(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BackupResourceVaultConfigsClient.Update`

```go
ctx := context.TODO()
id := backupresourcevaultconfigs.NewVaultID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultValue")

payload := backupresourcevaultconfigs.BackupResourceVaultConfigResource{
	// ...
}


read, err := client.Update(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

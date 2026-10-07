
## `github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2026-02-01/deletedvaults` Documentation

The `deletedvaults` SDK allows for interaction with Azure Resource Manager `keyvault` (API Version `2026-02-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2026-02-01/deletedvaults"
```


### Client Initialization

```go
client := deletedvaults.NewDeletedVaultsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DeletedVaultsClient.VaultsGetDeleted`

```go
ctx := context.TODO()
id := deletedvaults.NewDeletedVaultID("12345678-1234-9876-4563-123456789012", "locationName", "deletedVaultName")

read, err := client.VaultsGetDeleted(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeletedVaultsClient.VaultsPurgeDeleted`

```go
ctx := context.TODO()
id := deletedvaults.NewDeletedVaultID("12345678-1234-9876-4563-123456789012", "locationName", "deletedVaultName")

if err := client.VaultsPurgeDeletedThenPoll(ctx, id); err != nil {
	// handle the error
}
```

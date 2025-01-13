
## `github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/accountmigrations` Documentation

The `accountmigrations` SDK allows for interaction with Azure Resource Manager `storage` (API Version `2023-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/accountmigrations"
```


### Client Initialization

```go
client := accountmigrations.NewAccountMigrationsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AccountMigrationsClient.StorageAccountsCustomerInitiatedMigration`

```go
ctx := context.TODO()
id := commonids.NewStorageAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName")

payload := accountmigrations.StorageAccountMigration{
	// ...
}


if err := client.StorageAccountsCustomerInitiatedMigrationThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AccountMigrationsClient.StorageAccountsGetCustomerInitiatedMigration`

```go
ctx := context.TODO()
id := commonids.NewStorageAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName")

read, err := client.StorageAccountsGetCustomerInitiatedMigration(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


## `github.com/hashicorp/go-azure-sdk/resource-manager/storage/2022-05-01/deletedaccounts` Documentation

The `deletedaccounts` SDK allows for interaction with the Azure Resource Manager Service `storage` (API Version `2022-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/storage/2022-05-01/deletedaccounts"
```


### Client Initialization

```go
client := deletedaccounts.NewDeletedAccountsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DeletedAccountsClient.Get`

```go
ctx := context.TODO()
id := deletedaccounts.NewDeletedAccountID("12345678-1234-9876-4563-123456789012", "locationValue", "deletedAccountValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DeletedAccountsClient.List`

```go
ctx := context.TODO()
id := deletedaccounts.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

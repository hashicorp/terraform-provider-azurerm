
## `github.com/hashicorp/go-azure-sdk/resource-manager/videoindexer/2024-01-01/accounts` Documentation

The `accounts` SDK allows for interaction with Azure Resource Manager `videoindexer` (API Version `2024-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/videoindexer/2024-01-01/accounts"
```


### Client Initialization

```go
client := accounts.NewAccountsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AccountsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := accounts.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountName")

payload := accounts.Account{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AccountsClient.Delete`

```go
ctx := context.TODO()
id := accounts.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AccountsClient.Get`

```go
ctx := context.TODO()
id := accounts.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AccountsClient.List`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AccountsClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AccountsClient.Update`

```go
ctx := context.TODO()
id := accounts.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountName")

payload := accounts.AccountPatch{
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

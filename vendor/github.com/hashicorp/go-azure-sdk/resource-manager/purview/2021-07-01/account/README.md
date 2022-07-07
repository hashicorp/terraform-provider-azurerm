
## `github.com/hashicorp/go-azure-sdk/resource-manager/purview/2021-07-01/account` Documentation

The `account` SDK allows for interaction with the Azure Resource Manager Service `purview` (API Version `2021-07-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/purview/2021-07-01/account"
```


### Client Initialization

```go
client := account.NewAccountClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AccountClient.AddRootCollectionAdmin`

```go
ctx := context.TODO()
id := account.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue")

payload := account.CollectionAdminUpdate{
	// ...
}


read, err := client.AddRootCollectionAdmin(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AccountClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := account.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue")

payload := account.Account{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AccountClient.Delete`

```go
ctx := context.TODO()
id := account.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AccountClient.Get`

```go
ctx := context.TODO()
id := account.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AccountClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := account.NewResourceGroupID()

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AccountClient.ListBySubscription`

```go
ctx := context.TODO()
id := account.NewSubscriptionID()

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AccountClient.ListKeys`

```go
ctx := context.TODO()
id := account.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue")

read, err := client.ListKeys(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AccountClient.Update`

```go
ctx := context.TODO()
id := account.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue")

payload := account.AccountUpdateParameters{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


## `github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/account` Documentation

The `account` SDK allows for interaction with the Azure Resource Manager Service `datashare` (API Version `2019-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/account"
```


### Client Initialization

```go
client := account.NewAccountClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AccountClient.Create`

```go
ctx := context.TODO()
id := account.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue")

payload := account.Account{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
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
id := account.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

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
id := account.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AccountClient.Update`

```go
ctx := context.TODO()
id := account.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue")

payload := account.AccountUpdateParameters{
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

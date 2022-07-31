
## `github.com/hashicorp/go-azure-sdk/resource-manager/maps/2021-02-01/accounts` Documentation

The `accounts` SDK allows for interaction with the Azure Resource Manager Service `maps` (API Version `2021-02-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/maps/2021-02-01/accounts"
```


### Client Initialization

```go
client := accounts.NewAccountsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AccountsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := accounts.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue")

payload := accounts.MapsAccount{
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
id := accounts.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue")

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
id := accounts.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AccountsClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := accounts.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AccountsClient.ListBySubscription`

```go
ctx := context.TODO()
id := accounts.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AccountsClient.ListKeys`

```go
ctx := context.TODO()
id := accounts.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue")

read, err := client.ListKeys(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AccountsClient.RegenerateKeys`

```go
ctx := context.TODO()
id := accounts.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue")

payload := accounts.MapsKeySpecification{
	// ...
}


read, err := client.RegenerateKeys(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AccountsClient.Update`

```go
ctx := context.TODO()
id := accounts.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue")

payload := accounts.MapsAccountUpdateParameters{
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

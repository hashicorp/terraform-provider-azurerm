
## `github.com/hashicorp/go-azure-sdk/resource-manager/codesigning/2024-09-30-preview/codesigningaccounts` Documentation

The `codesigningaccounts` SDK allows for interaction with Azure Resource Manager `codesigning` (API Version `2024-09-30-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/codesigning/2024-09-30-preview/codesigningaccounts"
```


### Client Initialization

```go
client := codesigningaccounts.NewCodeSigningAccountsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `CodeSigningAccountsClient.CheckNameAvailability`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

payload := codesigningaccounts.CheckNameAvailability{
	// ...
}


read, err := client.CheckNameAvailability(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CodeSigningAccountsClient.Create`

```go
ctx := context.TODO()
id := codesigningaccounts.NewCodeSigningAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "codeSigningAccountName")

payload := codesigningaccounts.CodeSigningAccount{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CodeSigningAccountsClient.Delete`

```go
ctx := context.TODO()
id := codesigningaccounts.NewCodeSigningAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "codeSigningAccountName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CodeSigningAccountsClient.Get`

```go
ctx := context.TODO()
id := codesigningaccounts.NewCodeSigningAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "codeSigningAccountName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CodeSigningAccountsClient.ListByResourceGroup`

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


### Example Usage: `CodeSigningAccountsClient.ListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `CodeSigningAccountsClient.Update`

```go
ctx := context.TODO()
id := codesigningaccounts.NewCodeSigningAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "codeSigningAccountName")

payload := codesigningaccounts.CodeSigningAccountPatch{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

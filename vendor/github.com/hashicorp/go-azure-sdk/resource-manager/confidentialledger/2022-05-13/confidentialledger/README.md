
## `github.com/hashicorp/go-azure-sdk/resource-manager/confidentialledger/2022-05-13/confidentialledger` Documentation

The `confidentialledger` SDK allows for interaction with Azure Resource Manager `confidentialledger` (API Version `2022-05-13`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/confidentialledger/2022-05-13/confidentialledger"
```


### Client Initialization

```go
client := confidentialledger.NewConfidentialLedgerClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ConfidentialLedgerClient.LedgerCreate`

```go
ctx := context.TODO()
id := confidentialledger.NewLedgerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "ledgerName")

payload := confidentialledger.ConfidentialLedger{
	// ...
}


if err := client.LedgerCreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ConfidentialLedgerClient.LedgerDelete`

```go
ctx := context.TODO()
id := confidentialledger.NewLedgerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "ledgerName")

if err := client.LedgerDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ConfidentialLedgerClient.LedgerGet`

```go
ctx := context.TODO()
id := confidentialledger.NewLedgerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "ledgerName")

read, err := client.LedgerGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConfidentialLedgerClient.LedgerListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.LedgerListByResourceGroup(ctx, id, confidentialledger.DefaultLedgerListByResourceGroupOperationOptions())` can be used to do batched pagination
items, err := client.LedgerListByResourceGroupComplete(ctx, id, confidentialledger.DefaultLedgerListByResourceGroupOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ConfidentialLedgerClient.LedgerListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.LedgerListBySubscription(ctx, id, confidentialledger.DefaultLedgerListBySubscriptionOperationOptions())` can be used to do batched pagination
items, err := client.LedgerListBySubscriptionComplete(ctx, id, confidentialledger.DefaultLedgerListBySubscriptionOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ConfidentialLedgerClient.LedgerUpdate`

```go
ctx := context.TODO()
id := confidentialledger.NewLedgerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "ledgerName")

payload := confidentialledger.ConfidentialLedger{
	// ...
}


if err := client.LedgerUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

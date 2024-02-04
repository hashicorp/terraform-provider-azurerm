
## `github.com/hashicorp/go-azure-sdk/resource-manager/billing/2020-05-01/billingaccounts` Documentation

The `billingaccounts` SDK allows for interaction with the Azure Resource Manager Service `billing` (API Version `2020-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/billing/2020-05-01/billingaccounts"
```


### Client Initialization

```go
client := billingaccounts.NewBillingAccountsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `BillingAccountsClient.Get`

```go
ctx := context.TODO()
id := billingaccounts.NewBillingAccountID("billingAccountValue")

read, err := client.Get(ctx, id, billingaccounts.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BillingAccountsClient.List`

```go
ctx := context.TODO()


// alternatively `client.List(ctx, billingaccounts.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, billingaccounts.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `BillingAccountsClient.ListInvoiceSectionsByCreateSubscriptionPermission`

```go
ctx := context.TODO()
id := billingaccounts.NewBillingAccountID("billingAccountValue")

// alternatively `client.ListInvoiceSectionsByCreateSubscriptionPermission(ctx, id)` can be used to do batched pagination
items, err := client.ListInvoiceSectionsByCreateSubscriptionPermissionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `BillingAccountsClient.Update`

```go
ctx := context.TODO()
id := billingaccounts.NewBillingAccountID("billingAccountValue")

payload := billingaccounts.BillingAccountUpdateRequest{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

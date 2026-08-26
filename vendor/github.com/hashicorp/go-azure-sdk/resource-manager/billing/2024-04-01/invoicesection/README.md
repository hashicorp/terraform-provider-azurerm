
## `github.com/hashicorp/go-azure-sdk/resource-manager/billing/2024-04-01/invoicesection` Documentation

The `invoicesection` SDK allows for interaction with Azure Resource Manager `billing` (API Version `2024-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/billing/2024-04-01/invoicesection"
```


### Client Initialization

```go
client := invoicesection.NewInvoiceSectionClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `InvoiceSectionClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := invoicesection.NewInvoiceSectionID("billingAccountName", "billingProfileName", "invoiceSectionName")

payload := invoicesection.InvoiceSection{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `InvoiceSectionClient.Delete`

```go
ctx := context.TODO()
id := invoicesection.NewInvoiceSectionID("billingAccountName", "billingProfileName", "invoiceSectionName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `InvoiceSectionClient.Get`

```go
ctx := context.TODO()
id := invoicesection.NewInvoiceSectionID("billingAccountName", "billingProfileName", "invoiceSectionName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `InvoiceSectionClient.ListByBillingProfile`

```go
ctx := context.TODO()
id := invoicesection.NewBillingProfileID("billingAccountName", "billingProfileName")

// alternatively `client.ListByBillingProfile(ctx, id, invoicesection.DefaultListByBillingProfileOperationOptions())` can be used to do batched pagination
items, err := client.ListByBillingProfileComplete(ctx, id, invoicesection.DefaultListByBillingProfileOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `InvoiceSectionClient.ValidateDeleteEligibility`

```go
ctx := context.TODO()
id := invoicesection.NewInvoiceSectionID("billingAccountName", "billingProfileName", "invoiceSectionName")

read, err := client.ValidateDeleteEligibility(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

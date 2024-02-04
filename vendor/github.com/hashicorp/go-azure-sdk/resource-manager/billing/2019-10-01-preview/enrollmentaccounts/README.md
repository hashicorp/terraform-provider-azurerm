
## `github.com/hashicorp/go-azure-sdk/resource-manager/billing/2019-10-01-preview/enrollmentaccounts` Documentation

The `enrollmentaccounts` SDK allows for interaction with the Azure Resource Manager Service `billing` (API Version `2019-10-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/billing/2019-10-01-preview/enrollmentaccounts"
```


### Client Initialization

```go
client := enrollmentaccounts.NewEnrollmentAccountsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `EnrollmentAccountsClient.GetByEnrollmentAccountId`

```go
ctx := context.TODO()
id := enrollmentaccounts.NewEnrollmentAccountID("billingAccountValue", "enrollmentAccountValue")

read, err := client.GetByEnrollmentAccountId(ctx, id, enrollmentaccounts.DefaultGetByEnrollmentAccountIdOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EnrollmentAccountsClient.ListByBillingAccountName`

```go
ctx := context.TODO()
id := enrollmentaccounts.NewBillingAccountID("billingAccountValue")

read, err := client.ListByBillingAccountName(ctx, id, enrollmentaccounts.DefaultListByBillingAccountNameOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

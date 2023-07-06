
## `github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2023-05-01/cognitiveservicesaccounts` Documentation

The `cognitiveservicesaccounts` SDK allows for interaction with the Azure Resource Manager Service `cognitive` (API Version `2023-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2023-05-01/cognitiveservicesaccounts"
```


### Client Initialization

```go
client := cognitiveservicesaccounts.NewCognitiveServicesAccountsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `CognitiveServicesAccountsClient.AccountsCreate`

```go
ctx := context.TODO()
id := cognitiveservicesaccounts.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue")

payload := cognitiveservicesaccounts.Account{
	// ...
}


if err := client.AccountsCreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CognitiveServicesAccountsClient.AccountsDelete`

```go
ctx := context.TODO()
id := cognitiveservicesaccounts.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue")

if err := client.AccountsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CognitiveServicesAccountsClient.AccountsGet`

```go
ctx := context.TODO()
id := cognitiveservicesaccounts.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue")

read, err := client.AccountsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CognitiveServicesAccountsClient.AccountsList`

```go
ctx := context.TODO()
id := cognitiveservicesaccounts.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.AccountsList(ctx, id)` can be used to do batched pagination
items, err := client.AccountsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `CognitiveServicesAccountsClient.AccountsListByResourceGroup`

```go
ctx := context.TODO()
id := cognitiveservicesaccounts.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.AccountsListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.AccountsListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `CognitiveServicesAccountsClient.AccountsListKeys`

```go
ctx := context.TODO()
id := cognitiveservicesaccounts.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue")

read, err := client.AccountsListKeys(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CognitiveServicesAccountsClient.AccountsListModels`

```go
ctx := context.TODO()
id := cognitiveservicesaccounts.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue")

// alternatively `client.AccountsListModels(ctx, id)` can be used to do batched pagination
items, err := client.AccountsListModelsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `CognitiveServicesAccountsClient.AccountsListSkus`

```go
ctx := context.TODO()
id := cognitiveservicesaccounts.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue")

read, err := client.AccountsListSkus(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CognitiveServicesAccountsClient.AccountsListUsages`

```go
ctx := context.TODO()
id := cognitiveservicesaccounts.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue")

read, err := client.AccountsListUsages(ctx, id, cognitiveservicesaccounts.DefaultAccountsListUsagesOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CognitiveServicesAccountsClient.AccountsRegenerateKey`

```go
ctx := context.TODO()
id := cognitiveservicesaccounts.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue")

payload := cognitiveservicesaccounts.RegenerateKeyParameters{
	// ...
}


read, err := client.AccountsRegenerateKey(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CognitiveServicesAccountsClient.AccountsUpdate`

```go
ctx := context.TODO()
id := cognitiveservicesaccounts.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountValue")

payload := cognitiveservicesaccounts.Account{
	// ...
}


if err := client.AccountsUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CognitiveServicesAccountsClient.CheckDomainAvailability`

```go
ctx := context.TODO()
id := cognitiveservicesaccounts.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

payload := cognitiveservicesaccounts.CheckDomainAvailabilityParameter{
	// ...
}


read, err := client.CheckDomainAvailability(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CognitiveServicesAccountsClient.CheckSkuAvailability`

```go
ctx := context.TODO()
id := cognitiveservicesaccounts.NewLocationID("12345678-1234-9876-4563-123456789012", "locationValue")

payload := cognitiveservicesaccounts.CheckSkuAvailabilityParameter{
	// ...
}


read, err := client.CheckSkuAvailability(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CognitiveServicesAccountsClient.DeletedAccountsGet`

```go
ctx := context.TODO()
id := cognitiveservicesaccounts.NewDeletedAccountID("12345678-1234-9876-4563-123456789012", "locationValue", "example-resource-group", "deletedAccountValue")

read, err := client.DeletedAccountsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CognitiveServicesAccountsClient.DeletedAccountsList`

```go
ctx := context.TODO()
id := cognitiveservicesaccounts.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.DeletedAccountsList(ctx, id)` can be used to do batched pagination
items, err := client.DeletedAccountsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `CognitiveServicesAccountsClient.DeletedAccountsPurge`

```go
ctx := context.TODO()
id := cognitiveservicesaccounts.NewDeletedAccountID("12345678-1234-9876-4563-123456789012", "locationValue", "example-resource-group", "deletedAccountValue")

if err := client.DeletedAccountsPurgeThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CognitiveServicesAccountsClient.ResourceSkusList`

```go
ctx := context.TODO()
id := cognitiveservicesaccounts.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ResourceSkusList(ctx, id)` can be used to do batched pagination
items, err := client.ResourceSkusListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

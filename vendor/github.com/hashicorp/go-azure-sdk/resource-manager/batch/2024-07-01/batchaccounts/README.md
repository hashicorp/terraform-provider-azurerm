
## `github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/batchaccounts` Documentation

The `batchaccounts` SDK allows for interaction with Azure Resource Manager `batch` (API Version `2024-07-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/batchaccounts"
```


### Client Initialization

```go
client := batchaccounts.NewBatchAccountsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `BatchAccountsClient.BatchAccountCreate`

```go
ctx := context.TODO()
id := batchaccounts.NewBatchAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "batchAccountName")

payload := batchaccounts.BatchAccountCreateParameters{
	// ...
}


if err := client.BatchAccountCreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `BatchAccountsClient.BatchAccountDelete`

```go
ctx := context.TODO()
id := batchaccounts.NewBatchAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "batchAccountName")

if err := client.BatchAccountDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `BatchAccountsClient.BatchAccountGet`

```go
ctx := context.TODO()
id := batchaccounts.NewBatchAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "batchAccountName")

read, err := client.BatchAccountGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BatchAccountsClient.BatchAccountGetKeys`

```go
ctx := context.TODO()
id := batchaccounts.NewBatchAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "batchAccountName")

read, err := client.BatchAccountGetKeys(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BatchAccountsClient.BatchAccountList`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.BatchAccountList(ctx, id)` can be used to do batched pagination
items, err := client.BatchAccountListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `BatchAccountsClient.BatchAccountListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.BatchAccountListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.BatchAccountListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `BatchAccountsClient.BatchAccountListOutboundNetworkDependenciesEndpoints`

```go
ctx := context.TODO()
id := batchaccounts.NewBatchAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "batchAccountName")

// alternatively `client.BatchAccountListOutboundNetworkDependenciesEndpoints(ctx, id)` can be used to do batched pagination
items, err := client.BatchAccountListOutboundNetworkDependenciesEndpointsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `BatchAccountsClient.BatchAccountRegenerateKey`

```go
ctx := context.TODO()
id := batchaccounts.NewBatchAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "batchAccountName")

payload := batchaccounts.BatchAccountRegenerateKeyParameters{
	// ...
}


read, err := client.BatchAccountRegenerateKey(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BatchAccountsClient.BatchAccountSynchronizeAutoStorageKeys`

```go
ctx := context.TODO()
id := batchaccounts.NewBatchAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "batchAccountName")

read, err := client.BatchAccountSynchronizeAutoStorageKeys(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BatchAccountsClient.BatchAccountUpdate`

```go
ctx := context.TODO()
id := batchaccounts.NewBatchAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "batchAccountName")

payload := batchaccounts.BatchAccountUpdateParameters{
	// ...
}


read, err := client.BatchAccountUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

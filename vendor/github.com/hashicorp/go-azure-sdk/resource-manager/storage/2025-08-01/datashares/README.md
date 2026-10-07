
## `github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/datashares` Documentation

The `datashares` SDK allows for interaction with Azure Resource Manager `storage` (API Version `2025-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/datashares"
```


### Client Initialization

```go
client := datashares.NewDataSharesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DataSharesClient.Create`

```go
ctx := context.TODO()
id := datashares.NewDataShareID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName", "dataShareName")

payload := datashares.DataShare{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DataSharesClient.Delete`

```go
ctx := context.TODO()
id := datashares.NewDataShareID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName", "dataShareName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `DataSharesClient.Get`

```go
ctx := context.TODO()
id := datashares.NewDataShareID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName", "dataShareName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DataSharesClient.ListByStorageAccount`

```go
ctx := context.TODO()
id := commonids.NewStorageAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName")

// alternatively `client.ListByStorageAccount(ctx, id)` can be used to do batched pagination
items, err := client.ListByStorageAccountComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `DataSharesClient.Update`

```go
ctx := context.TODO()
id := datashares.NewDataShareID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName", "dataShareName")

payload := datashares.DataShareUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

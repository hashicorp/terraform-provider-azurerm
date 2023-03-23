
## `github.com/hashicorp/go-azure-sdk/resource-manager/storage/2022-05-01/tableservice` Documentation

The `tableservice` SDK allows for interaction with the Azure Resource Manager Service `storage` (API Version `2022-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/storage/2022-05-01/tableservice"
```


### Client Initialization

```go
client := tableservice.NewTableServiceClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `TableServiceClient.TableCreate`

```go
ctx := context.TODO()
id := tableservice.NewTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountValue", "tableValue")

payload := tableservice.Table{
	// ...
}


read, err := client.TableCreate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TableServiceClient.TableDelete`

```go
ctx := context.TODO()
id := tableservice.NewTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountValue", "tableValue")

read, err := client.TableDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TableServiceClient.TableGet`

```go
ctx := context.TODO()
id := tableservice.NewTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountValue", "tableValue")

read, err := client.TableGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TableServiceClient.TableList`

```go
ctx := context.TODO()
id := tableservice.NewStorageAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountValue")

// alternatively `client.TableList(ctx, id)` can be used to do batched pagination
items, err := client.TableListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `TableServiceClient.TableUpdate`

```go
ctx := context.TODO()
id := tableservice.NewTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountValue", "tableValue")

payload := tableservice.Table{
	// ...
}


read, err := client.TableUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

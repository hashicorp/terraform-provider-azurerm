
## `github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/inventoryitems` Documentation

The `inventoryitems` SDK allows for interaction with Azure Resource Manager `systemcentervirtualmachinemanager` (API Version `2023-10-07`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/inventoryitems"
```


### Client Initialization

```go
client := inventoryitems.NewInventoryItemsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `InventoryItemsClient.Create`

```go
ctx := context.TODO()
id := inventoryitems.NewInventoryItemID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vmmServerName", "inventoryItemName")

payload := inventoryitems.InventoryItem{
	// ...
}


read, err := client.Create(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `InventoryItemsClient.Delete`

```go
ctx := context.TODO()
id := inventoryitems.NewInventoryItemID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vmmServerName", "inventoryItemName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `InventoryItemsClient.Get`

```go
ctx := context.TODO()
id := inventoryitems.NewInventoryItemID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vmmServerName", "inventoryItemName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `InventoryItemsClient.ListByVMmServer`

```go
ctx := context.TODO()
id := inventoryitems.NewVMmServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vmmServerName")

// alternatively `client.ListByVMmServer(ctx, id)` can be used to do batched pagination
items, err := client.ListByVMmServerComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

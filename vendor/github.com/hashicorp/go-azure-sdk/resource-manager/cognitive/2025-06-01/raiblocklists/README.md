
## `github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2025-06-01/raiblocklists` Documentation

The `raiblocklists` SDK allows for interaction with Azure Resource Manager `cognitive` (API Version `2025-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2025-06-01/raiblocklists"
```


### Client Initialization

```go
client := raiblocklists.NewRaiBlocklistsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `RaiBlocklistsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := raiblocklists.NewRaiBlocklistID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountName", "raiBlocklistName")

payload := raiblocklists.RaiBlocklist{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RaiBlocklistsClient.Delete`

```go
ctx := context.TODO()
id := raiblocklists.NewRaiBlocklistID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountName", "raiBlocklistName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `RaiBlocklistsClient.Get`

```go
ctx := context.TODO()
id := raiblocklists.NewRaiBlocklistID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountName", "raiBlocklistName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RaiBlocklistsClient.List`

```go
ctx := context.TODO()
id := raiblocklists.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `RaiBlocklistsClient.RaiBlocklistItemsBatchAdd`

```go
ctx := context.TODO()
id := raiblocklists.NewRaiBlocklistID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountName", "raiBlocklistName")
var payload []RaiBlocklistItemBulkRequest

read, err := client.RaiBlocklistItemsBatchAdd(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RaiBlocklistsClient.RaiBlocklistItemsBatchDelete`

```go
ctx := context.TODO()
id := raiblocklists.NewRaiBlocklistID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountName", "raiBlocklistName")
var payload interface{}

read, err := client.RaiBlocklistItemsBatchDelete(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RaiBlocklistsClient.RaiBlocklistItemsCreateOrUpdate`

```go
ctx := context.TODO()
id := raiblocklists.NewRaiBlocklistItemID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountName", "raiBlocklistName", "raiBlocklistItemName")

payload := raiblocklists.RaiBlocklistItem{
	// ...
}


read, err := client.RaiBlocklistItemsCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RaiBlocklistsClient.RaiBlocklistItemsDelete`

```go
ctx := context.TODO()
id := raiblocklists.NewRaiBlocklistItemID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountName", "raiBlocklistName", "raiBlocklistItemName")

if err := client.RaiBlocklistItemsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `RaiBlocklistsClient.RaiBlocklistItemsGet`

```go
ctx := context.TODO()
id := raiblocklists.NewRaiBlocklistItemID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountName", "raiBlocklistName", "raiBlocklistItemName")

read, err := client.RaiBlocklistItemsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RaiBlocklistsClient.RaiBlocklistItemsList`

```go
ctx := context.TODO()
id := raiblocklists.NewRaiBlocklistID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountName", "raiBlocklistName")

// alternatively `client.RaiBlocklistItemsList(ctx, id)` can be used to do batched pagination
items, err := client.RaiBlocklistItemsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

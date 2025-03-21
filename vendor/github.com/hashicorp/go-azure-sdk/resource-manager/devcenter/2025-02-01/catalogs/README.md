
## `github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/catalogs` Documentation

The `catalogs` SDK allows for interaction with Azure Resource Manager `devcenter` (API Version `2025-02-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/catalogs"
```


### Client Initialization

```go
client := catalogs.NewCatalogsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `CatalogsClient.Connect`

```go
ctx := context.TODO()
id := catalogs.NewDevCenterCatalogID("12345678-1234-9876-4563-123456789012", "example-resource-group", "devCenterName", "catalogName")

if err := client.ConnectThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CatalogsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := catalogs.NewDevCenterCatalogID("12345678-1234-9876-4563-123456789012", "example-resource-group", "devCenterName", "catalogName")

payload := catalogs.Catalog{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CatalogsClient.Delete`

```go
ctx := context.TODO()
id := catalogs.NewDevCenterCatalogID("12345678-1234-9876-4563-123456789012", "example-resource-group", "devCenterName", "catalogName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CatalogsClient.Get`

```go
ctx := context.TODO()
id := catalogs.NewDevCenterCatalogID("12345678-1234-9876-4563-123456789012", "example-resource-group", "devCenterName", "catalogName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CatalogsClient.GetSyncErrorDetails`

```go
ctx := context.TODO()
id := catalogs.NewDevCenterCatalogID("12345678-1234-9876-4563-123456789012", "example-resource-group", "devCenterName", "catalogName")

read, err := client.GetSyncErrorDetails(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CatalogsClient.ListByDevCenter`

```go
ctx := context.TODO()
id := catalogs.NewDevCenterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "devCenterName")

// alternatively `client.ListByDevCenter(ctx, id, catalogs.DefaultListByDevCenterOperationOptions())` can be used to do batched pagination
items, err := client.ListByDevCenterComplete(ctx, id, catalogs.DefaultListByDevCenterOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `CatalogsClient.Sync`

```go
ctx := context.TODO()
id := catalogs.NewDevCenterCatalogID("12345678-1234-9876-4563-123456789012", "example-resource-group", "devCenterName", "catalogName")

if err := client.SyncThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CatalogsClient.Update`

```go
ctx := context.TODO()
id := catalogs.NewDevCenterCatalogID("12345678-1234-9876-4563-123456789012", "example-resource-group", "devCenterName", "catalogName")

payload := catalogs.CatalogUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

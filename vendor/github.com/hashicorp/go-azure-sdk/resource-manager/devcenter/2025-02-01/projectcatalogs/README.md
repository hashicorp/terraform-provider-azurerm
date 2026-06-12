
## `github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/projectcatalogs` Documentation

The `projectcatalogs` SDK allows for interaction with Azure Resource Manager `devcenter` (API Version `2025-02-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/projectcatalogs"
```


### Client Initialization

```go
client := projectcatalogs.NewProjectCatalogsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ProjectCatalogsClient.Connect`

```go
ctx := context.TODO()
id := projectcatalogs.NewCatalogID("12345678-1234-9876-4563-123456789012", "example-resource-group", "projectName", "catalogName")

if err := client.ConnectThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ProjectCatalogsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := projectcatalogs.NewCatalogID("12345678-1234-9876-4563-123456789012", "example-resource-group", "projectName", "catalogName")

payload := projectcatalogs.Catalog{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ProjectCatalogsClient.Delete`

```go
ctx := context.TODO()
id := projectcatalogs.NewCatalogID("12345678-1234-9876-4563-123456789012", "example-resource-group", "projectName", "catalogName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ProjectCatalogsClient.Get`

```go
ctx := context.TODO()
id := projectcatalogs.NewCatalogID("12345678-1234-9876-4563-123456789012", "example-resource-group", "projectName", "catalogName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProjectCatalogsClient.GetSyncErrorDetails`

```go
ctx := context.TODO()
id := projectcatalogs.NewCatalogID("12345678-1234-9876-4563-123456789012", "example-resource-group", "projectName", "catalogName")

read, err := client.GetSyncErrorDetails(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProjectCatalogsClient.List`

```go
ctx := context.TODO()
id := projectcatalogs.NewProjectID("12345678-1234-9876-4563-123456789012", "example-resource-group", "projectName")

// alternatively `client.List(ctx, id, projectcatalogs.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, projectcatalogs.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ProjectCatalogsClient.Patch`

```go
ctx := context.TODO()
id := projectcatalogs.NewCatalogID("12345678-1234-9876-4563-123456789012", "example-resource-group", "projectName", "catalogName")

payload := projectcatalogs.CatalogUpdate{
	// ...
}


if err := client.PatchThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ProjectCatalogsClient.Sync`

```go
ctx := context.TODO()
id := projectcatalogs.NewCatalogID("12345678-1234-9876-4563-123456789012", "example-resource-group", "projectName", "catalogName")

if err := client.SyncThenPoll(ctx, id); err != nil {
	// handle the error
}
```

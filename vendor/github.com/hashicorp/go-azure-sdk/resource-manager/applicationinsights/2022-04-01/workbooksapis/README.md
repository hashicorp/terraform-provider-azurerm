
## `github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2022-04-01/workbooksapis` Documentation

The `workbooksapis` SDK allows for interaction with the Azure Resource Manager Service `applicationinsights` (API Version `2022-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2022-04-01/workbooksapis"
```


### Client Initialization

```go
client := workbooksapis.NewWorkbooksAPIsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `WorkbooksAPIsClient.WorkbooksCreateOrUpdate`

```go
ctx := context.TODO()
id := workbooksapis.NewWorkbookID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workbookValue")

payload := workbooksapis.Workbook{
	// ...
}


read, err := client.WorkbooksCreateOrUpdate(ctx, id, payload, workbooksapis.DefaultWorkbooksCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WorkbooksAPIsClient.WorkbooksDelete`

```go
ctx := context.TODO()
id := workbooksapis.NewWorkbookID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workbookValue")

read, err := client.WorkbooksDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WorkbooksAPIsClient.WorkbooksGet`

```go
ctx := context.TODO()
id := workbooksapis.NewWorkbookID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workbookValue")

read, err := client.WorkbooksGet(ctx, id, workbooksapis.DefaultWorkbooksGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WorkbooksAPIsClient.WorkbooksListByResourceGroup`

```go
ctx := context.TODO()
id := workbooksapis.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.WorkbooksListByResourceGroup(ctx, id, workbooksapis.DefaultWorkbooksListByResourceGroupOperationOptions())` can be used to do batched pagination
items, err := client.WorkbooksListByResourceGroupComplete(ctx, id, workbooksapis.DefaultWorkbooksListByResourceGroupOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WorkbooksAPIsClient.WorkbooksListBySubscription`

```go
ctx := context.TODO()
id := workbooksapis.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.WorkbooksListBySubscription(ctx, id, workbooksapis.DefaultWorkbooksListBySubscriptionOperationOptions())` can be used to do batched pagination
items, err := client.WorkbooksListBySubscriptionComplete(ctx, id, workbooksapis.DefaultWorkbooksListBySubscriptionOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WorkbooksAPIsClient.WorkbooksRevisionGet`

```go
ctx := context.TODO()
id := workbooksapis.NewRevisionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workbookValue", "revisionIdValue")

read, err := client.WorkbooksRevisionGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WorkbooksAPIsClient.WorkbooksRevisionsList`

```go
ctx := context.TODO()
id := workbooksapis.NewWorkbookID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workbookValue")

// alternatively `client.WorkbooksRevisionsList(ctx, id)` can be used to do batched pagination
items, err := client.WorkbooksRevisionsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WorkbooksAPIsClient.WorkbooksUpdate`

```go
ctx := context.TODO()
id := workbooksapis.NewWorkbookID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workbookValue")

payload := workbooksapis.WorkbookUpdateParameters{
	// ...
}


read, err := client.WorkbooksUpdate(ctx, id, payload, workbooksapis.DefaultWorkbooksUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

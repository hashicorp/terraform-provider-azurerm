
## `github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2022-04-01/applicationinsights` Documentation

The `applicationinsights` SDK allows for interaction with the Azure Resource Manager Service `applicationinsights` (API Version `2022-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2022-04-01/applicationinsights"
```


### Client Initialization

```go
client := applicationinsights.NewApplicationInsightsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ApplicationInsightsClient.WorkbooksCreateOrUpdate`

```go
ctx := context.TODO()
id := applicationinsights.NewWorkbookID("12345678-1234-9876-4563-123456789012", "example-resource-group", "resourceValue")

payload := applicationinsights.Workbook{
	// ...
}


read, err := client.WorkbooksCreateOrUpdate(ctx, id, payload, applicationinsights.DefaultWorkbooksCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApplicationInsightsClient.WorkbooksDelete`

```go
ctx := context.TODO()
id := applicationinsights.NewWorkbookID("12345678-1234-9876-4563-123456789012", "example-resource-group", "resourceValue")

read, err := client.WorkbooksDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApplicationInsightsClient.WorkbooksGet`

```go
ctx := context.TODO()
id := applicationinsights.NewWorkbookID("12345678-1234-9876-4563-123456789012", "example-resource-group", "resourceValue")

read, err := client.WorkbooksGet(ctx, id, applicationinsights.DefaultWorkbooksGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApplicationInsightsClient.WorkbooksListByResourceGroup`

```go
ctx := context.TODO()
id := applicationinsights.NewResourceGroupID()

// alternatively `client.WorkbooksListByResourceGroup(ctx, id, applicationinsights.DefaultWorkbooksListByResourceGroupOperationOptions())` can be used to do batched pagination
items, err := client.WorkbooksListByResourceGroupComplete(ctx, id, applicationinsights.DefaultWorkbooksListByResourceGroupOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ApplicationInsightsClient.WorkbooksListBySubscription`

```go
ctx := context.TODO()
id := applicationinsights.NewSubscriptionID()

// alternatively `client.WorkbooksListBySubscription(ctx, id, applicationinsights.DefaultWorkbooksListBySubscriptionOperationOptions())` can be used to do batched pagination
items, err := client.WorkbooksListBySubscriptionComplete(ctx, id, applicationinsights.DefaultWorkbooksListBySubscriptionOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ApplicationInsightsClient.WorkbooksRevisionGet`

```go
ctx := context.TODO()
id := applicationinsights.NewRevisionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "resourceValue", "revisionIdValue")

read, err := client.WorkbooksRevisionGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApplicationInsightsClient.WorkbooksRevisionsList`

```go
ctx := context.TODO()
id := applicationinsights.NewWorkbookID("12345678-1234-9876-4563-123456789012", "example-resource-group", "resourceValue")

// alternatively `client.WorkbooksRevisionsList(ctx, id)` can be used to do batched pagination
items, err := client.WorkbooksRevisionsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ApplicationInsightsClient.WorkbooksUpdate`

```go
ctx := context.TODO()
id := applicationinsights.NewWorkbookID("12345678-1234-9876-4563-123456789012", "example-resource-group", "resourceValue")

payload := applicationinsights.WorkbookUpdateParameters{
	// ...
}


read, err := client.WorkbooksUpdate(ctx, id, payload, applicationinsights.DefaultWorkbooksUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

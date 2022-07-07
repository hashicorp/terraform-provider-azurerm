
## `github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-11-20/applicationinsights` Documentation

The `applicationinsights` SDK allows for interaction with the Azure Resource Manager Service `applicationinsights` (API Version `2020-11-20`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-11-20/applicationinsights"
```


### Client Initialization

```go
client := applicationinsights.NewApplicationInsightsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ApplicationInsightsClient.WorkbookTemplatesCreateOrUpdate`

```go
ctx := context.TODO()
id := applicationinsights.NewWorkbookTemplateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "resourceValue")

payload := applicationinsights.WorkbookTemplate{
	// ...
}


read, err := client.WorkbookTemplatesCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApplicationInsightsClient.WorkbookTemplatesDelete`

```go
ctx := context.TODO()
id := applicationinsights.NewWorkbookTemplateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "resourceValue")

read, err := client.WorkbookTemplatesDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApplicationInsightsClient.WorkbookTemplatesGet`

```go
ctx := context.TODO()
id := applicationinsights.NewWorkbookTemplateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "resourceValue")

read, err := client.WorkbookTemplatesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApplicationInsightsClient.WorkbookTemplatesListByResourceGroup`

```go
ctx := context.TODO()
id := applicationinsights.NewResourceGroupID()

read, err := client.WorkbookTemplatesListByResourceGroup(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApplicationInsightsClient.WorkbookTemplatesUpdate`

```go
ctx := context.TODO()
id := applicationinsights.NewWorkbookTemplateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "resourceValue")

payload := applicationinsights.WorkbookTemplateUpdateParameters{
	// ...
}


read, err := client.WorkbookTemplatesUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


## `github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-11-20/workbooktemplatesapis` Documentation

The `workbooktemplatesapis` SDK allows for interaction with the Azure Resource Manager Service `applicationinsights` (API Version `2020-11-20`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-11-20/workbooktemplatesapis"
```


### Client Initialization

```go
client := workbooktemplatesapis.NewWorkbookTemplatesAPIsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `WorkbookTemplatesAPIsClient.WorkbookTemplatesCreateOrUpdate`

```go
ctx := context.TODO()
id := workbooktemplatesapis.NewWorkbookTemplateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workbookTemplateValue")

payload := workbooktemplatesapis.WorkbookTemplate{
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


### Example Usage: `WorkbookTemplatesAPIsClient.WorkbookTemplatesDelete`

```go
ctx := context.TODO()
id := workbooktemplatesapis.NewWorkbookTemplateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workbookTemplateValue")

read, err := client.WorkbookTemplatesDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WorkbookTemplatesAPIsClient.WorkbookTemplatesGet`

```go
ctx := context.TODO()
id := workbooktemplatesapis.NewWorkbookTemplateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workbookTemplateValue")

read, err := client.WorkbookTemplatesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WorkbookTemplatesAPIsClient.WorkbookTemplatesListByResourceGroup`

```go
ctx := context.TODO()
id := workbooktemplatesapis.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

read, err := client.WorkbookTemplatesListByResourceGroup(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WorkbookTemplatesAPIsClient.WorkbookTemplatesUpdate`

```go
ctx := context.TODO()
id := workbooktemplatesapis.NewWorkbookTemplateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workbookTemplateValue")

payload := workbooktemplatesapis.WorkbookTemplateUpdateParameters{
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

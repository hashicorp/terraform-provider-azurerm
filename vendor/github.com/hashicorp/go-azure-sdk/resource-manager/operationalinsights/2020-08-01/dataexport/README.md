
## `github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/dataexport` Documentation

The `dataexport` SDK allows for interaction with the Azure Resource Manager Service `operationalinsights` (API Version `2020-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/dataexport"
```


### Client Initialization

```go
client := dataexport.NewDataExportClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DataExportClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := dataexport.NewDataExportID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "dataExportValue")

payload := dataexport.DataExport{
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


### Example Usage: `DataExportClient.Delete`

```go
ctx := context.TODO()
id := dataexport.NewDataExportID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "dataExportValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DataExportClient.Get`

```go
ctx := context.TODO()
id := dataexport.NewDataExportID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "dataExportValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DataExportClient.ListByWorkspace`

```go
ctx := context.TODO()
id := dataexport.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue")

read, err := client.ListByWorkspace(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

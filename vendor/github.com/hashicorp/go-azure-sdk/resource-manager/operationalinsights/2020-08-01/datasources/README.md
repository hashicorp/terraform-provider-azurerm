
## `github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/datasources` Documentation

The `datasources` SDK allows for interaction with the Azure Resource Manager Service `operationalinsights` (API Version `2020-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/datasources"
```


### Client Initialization

```go
client := datasources.NewDataSourcesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DataSourcesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := datasources.NewDataSourceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "dataSourceValue")

payload := datasources.DataSource{
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


### Example Usage: `DataSourcesClient.Delete`

```go
ctx := context.TODO()
id := datasources.NewDataSourceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "dataSourceValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DataSourcesClient.Get`

```go
ctx := context.TODO()
id := datasources.NewDataSourceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "dataSourceValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DataSourcesClient.ListByWorkspace`

```go
ctx := context.TODO()
id := datasources.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue")

// alternatively `client.ListByWorkspace(ctx, id, datasources.DefaultListByWorkspaceOperationOptions())` can be used to do batched pagination
items, err := client.ListByWorkspaceComplete(ctx, id, datasources.DefaultListByWorkspaceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

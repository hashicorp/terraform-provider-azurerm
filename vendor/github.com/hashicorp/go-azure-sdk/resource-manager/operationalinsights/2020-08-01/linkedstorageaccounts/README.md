
## `github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/linkedstorageaccounts` Documentation

The `linkedstorageaccounts` SDK allows for interaction with Azure Resource Manager `operationalinsights` (API Version `2020-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/linkedstorageaccounts"
```


### Client Initialization

```go
client := linkedstorageaccounts.NewLinkedStorageAccountsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `LinkedStorageAccountsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := linkedstorageaccounts.NewDataSourceTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "Alerts")

payload := linkedstorageaccounts.LinkedStorageAccountsResource{
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


### Example Usage: `LinkedStorageAccountsClient.Delete`

```go
ctx := context.TODO()
id := linkedstorageaccounts.NewDataSourceTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "Alerts")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LinkedStorageAccountsClient.Get`

```go
ctx := context.TODO()
id := linkedstorageaccounts.NewDataSourceTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "Alerts")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `LinkedStorageAccountsClient.ListByWorkspace`

```go
ctx := context.TODO()
id := linkedstorageaccounts.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName")

read, err := client.ListByWorkspace(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

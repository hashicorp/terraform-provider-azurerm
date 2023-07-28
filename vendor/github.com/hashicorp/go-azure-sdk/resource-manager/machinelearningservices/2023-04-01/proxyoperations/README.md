
## `github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2023-04-01/proxyoperations` Documentation

The `proxyoperations` SDK allows for interaction with the Azure Resource Manager Service `machinelearningservices` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2023-04-01/proxyoperations"
```


### Client Initialization

```go
client := proxyoperations.NewProxyOperationsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ProxyOperationsClient.WorkspacesListNotebookKeys`

```go
ctx := context.TODO()
id := proxyoperations.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue")

read, err := client.WorkspacesListNotebookKeys(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProxyOperationsClient.WorkspacesListStorageAccountKeys`

```go
ctx := context.TODO()
id := proxyoperations.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue")

read, err := client.WorkspacesListStorageAccountKeys(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProxyOperationsClient.WorkspacesPrepareNotebook`

```go
ctx := context.TODO()
id := proxyoperations.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue")

if err := client.WorkspacesPrepareNotebookThenPoll(ctx, id); err != nil {
	// handle the error
}
```

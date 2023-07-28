
## `github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2023-04-01/v2workspaceconnectionresource` Documentation

The `v2workspaceconnectionresource` SDK allows for interaction with the Azure Resource Manager Service `machinelearningservices` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2023-04-01/v2workspaceconnectionresource"
```


### Client Initialization

```go
client := v2workspaceconnectionresource.NewV2WorkspaceConnectionResourceClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `V2WorkspaceConnectionResourceClient.WorkspaceConnectionsCreate`

```go
ctx := context.TODO()
id := v2workspaceconnectionresource.NewConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "connectionValue")

payload := v2workspaceconnectionresource.WorkspaceConnectionPropertiesV2BasicResource{
	// ...
}


read, err := client.WorkspaceConnectionsCreate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `V2WorkspaceConnectionResourceClient.WorkspaceConnectionsDelete`

```go
ctx := context.TODO()
id := v2workspaceconnectionresource.NewConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "connectionValue")

read, err := client.WorkspaceConnectionsDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `V2WorkspaceConnectionResourceClient.WorkspaceConnectionsGet`

```go
ctx := context.TODO()
id := v2workspaceconnectionresource.NewConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "connectionValue")

read, err := client.WorkspaceConnectionsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `V2WorkspaceConnectionResourceClient.WorkspaceConnectionsList`

```go
ctx := context.TODO()
id := v2workspaceconnectionresource.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue")

// alternatively `client.WorkspaceConnectionsList(ctx, id, v2workspaceconnectionresource.DefaultWorkspaceConnectionsListOperationOptions())` can be used to do batched pagination
items, err := client.WorkspaceConnectionsListComplete(ctx, id, v2workspaceconnectionresource.DefaultWorkspaceConnectionsListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

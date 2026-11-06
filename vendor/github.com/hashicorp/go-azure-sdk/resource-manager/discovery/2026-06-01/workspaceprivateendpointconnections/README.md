
## `github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/workspaceprivateendpointconnections` Documentation

The `workspaceprivateendpointconnections` SDK allows for interaction with Azure Resource Manager `discovery` (API Version `2026-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/workspaceprivateendpointconnections"
```


### Client Initialization

```go
client := workspaceprivateendpointconnections.NewWorkspacePrivateEndpointConnectionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `WorkspacePrivateEndpointConnectionsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := workspaceprivateendpointconnections.NewPrivateEndpointConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "privateEndpointConnectionName")

payload := workspaceprivateendpointconnections.WorkspacePrivateEndpointConnection{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `WorkspacePrivateEndpointConnectionsClient.Delete`

```go
ctx := context.TODO()
id := workspaceprivateendpointconnections.NewPrivateEndpointConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "privateEndpointConnectionName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `WorkspacePrivateEndpointConnectionsClient.Get`

```go
ctx := context.TODO()
id := workspaceprivateendpointconnections.NewPrivateEndpointConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "privateEndpointConnectionName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WorkspacePrivateEndpointConnectionsClient.ListByWorkspace`

```go
ctx := context.TODO()
id := workspaceprivateendpointconnections.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName")

// alternatively `client.ListByWorkspace(ctx, id)` can be used to do batched pagination
items, err := client.ListByWorkspaceComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

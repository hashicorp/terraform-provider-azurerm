
## `github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2022-12-01/workspaces` Documentation

The `workspaces` SDK allows for interaction with the Azure Resource Manager Service `healthcareapis` (API Version `2022-12-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2022-12-01/workspaces"
```


### Client Initialization

```go
client := workspaces.NewWorkspacesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `WorkspacesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := workspaces.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue")

payload := workspaces.Workspace{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `WorkspacesClient.Delete`

```go
ctx := context.TODO()
id := workspaces.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `WorkspacesClient.Get`

```go
ctx := context.TODO()
id := workspaces.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `WorkspacesClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := workspaces.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WorkspacesClient.ListBySubscription`

```go
ctx := context.TODO()
id := workspaces.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `WorkspacesClient.Update`

```go
ctx := context.TODO()
id := workspaces.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue")

payload := workspaces.ResourceTags{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

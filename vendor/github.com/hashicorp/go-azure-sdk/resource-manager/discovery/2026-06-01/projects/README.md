
## `github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/projects` Documentation

The `projects` SDK allows for interaction with Azure Resource Manager `discovery` (API Version `2026-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/projects"
```


### Client Initialization

```go
client := projects.NewProjectsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ProjectsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := projects.NewProjectID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "projectName")

payload := projects.Project{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ProjectsClient.Delete`

```go
ctx := context.TODO()
id := projects.NewProjectID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "projectName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ProjectsClient.Get`

```go
ctx := context.TODO()
id := projects.NewProjectID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "projectName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProjectsClient.ListByWorkspace`

```go
ctx := context.TODO()
id := projects.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName")

// alternatively `client.ListByWorkspace(ctx, id)` can be used to do batched pagination
items, err := client.ListByWorkspaceComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ProjectsClient.Update`

```go
ctx := context.TODO()
id := projects.NewProjectID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "projectName")

payload := projects.ProjectUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

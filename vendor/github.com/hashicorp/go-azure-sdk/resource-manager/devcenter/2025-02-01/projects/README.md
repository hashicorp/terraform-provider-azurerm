
## `github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/projects` Documentation

The `projects` SDK allows for interaction with Azure Resource Manager `devcenter` (API Version `2025-02-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/projects"
```


### Client Initialization

```go
client := projects.NewProjectsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ProjectsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := projects.NewProjectID("12345678-1234-9876-4563-123456789012", "example-resource-group", "projectName")

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
id := projects.NewProjectID("12345678-1234-9876-4563-123456789012", "example-resource-group", "projectName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ProjectsClient.Get`

```go
ctx := context.TODO()
id := projects.NewProjectID("12345678-1234-9876-4563-123456789012", "example-resource-group", "projectName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProjectsClient.GetInheritedSettings`

```go
ctx := context.TODO()
id := projects.NewProjectID("12345678-1234-9876-4563-123456789012", "example-resource-group", "projectName")

read, err := client.GetInheritedSettings(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProjectsClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id, projects.DefaultListByResourceGroupOperationOptions())` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id, projects.DefaultListByResourceGroupOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ProjectsClient.ListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id, projects.DefaultListBySubscriptionOperationOptions())` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id, projects.DefaultListBySubscriptionOperationOptions())
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
id := projects.NewProjectID("12345678-1234-9876-4563-123456789012", "example-resource-group", "projectName")

payload := projects.ProjectUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

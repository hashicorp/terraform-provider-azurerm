
## `github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2023-03-01/projects` Documentation

The `projects` SDK allows for interaction with Azure Resource Manager `storagemover` (API Version `2023-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2023-03-01/projects"
```


### Client Initialization

```go
client := projects.NewProjectsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ProjectsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := projects.NewProjectID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageMoverName", "projectName")

payload := projects.Project{
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


### Example Usage: `ProjectsClient.Delete`

```go
ctx := context.TODO()
id := projects.NewProjectID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageMoverName", "projectName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ProjectsClient.Get`

```go
ctx := context.TODO()
id := projects.NewProjectID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageMoverName", "projectName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProjectsClient.List`

```go
ctx := context.TODO()
id := projects.NewStorageMoverID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageMoverName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
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
id := projects.NewProjectID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageMoverName", "projectName")

payload := projects.ProjectUpdateParameters{
	// ...
}


read, err := client.Update(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

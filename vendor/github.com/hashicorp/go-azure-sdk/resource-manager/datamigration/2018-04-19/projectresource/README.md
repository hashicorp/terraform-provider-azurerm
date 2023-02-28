
## `github.com/hashicorp/go-azure-sdk/resource-manager/datamigration/2018-04-19/projectresource` Documentation

The `projectresource` SDK allows for interaction with the Azure Resource Manager Service `datamigration` (API Version `2018-04-19`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/datamigration/2018-04-19/projectresource"
```


### Client Initialization

```go
client := projectresource.NewProjectResourceClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ProjectResourceClient.ProjectsCreateOrUpdate`

```go
ctx := context.TODO()
id := projectresource.NewProjectID("12345678-1234-9876-4563-123456789012", "resourceGroupValue", "serviceValue", "projectValue")

payload := projectresource.Project{
	// ...
}


read, err := client.ProjectsCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProjectResourceClient.ProjectsDelete`

```go
ctx := context.TODO()
id := projectresource.NewProjectID("12345678-1234-9876-4563-123456789012", "resourceGroupValue", "serviceValue", "projectValue")

read, err := client.ProjectsDelete(ctx, id, projectresource.DefaultProjectsDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProjectResourceClient.ProjectsGet`

```go
ctx := context.TODO()
id := projectresource.NewProjectID("12345678-1234-9876-4563-123456789012", "resourceGroupValue", "serviceValue", "projectValue")

read, err := client.ProjectsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProjectResourceClient.ProjectsListByResourceGroup`

```go
ctx := context.TODO()
id := projectresource.NewServiceID("12345678-1234-9876-4563-123456789012", "resourceGroupValue", "serviceValue")

// alternatively `client.ProjectsListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ProjectsListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ProjectResourceClient.ProjectsUpdate`

```go
ctx := context.TODO()
id := projectresource.NewProjectID("12345678-1234-9876-4563-123456789012", "resourceGroupValue", "serviceValue", "projectValue")

payload := projectresource.Project{
	// ...
}


read, err := client.ProjectsUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

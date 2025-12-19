
## `github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2025-06-01/projectconnectionresource` Documentation

The `projectconnectionresource` SDK allows for interaction with Azure Resource Manager `cognitive` (API Version `2025-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2025-06-01/projectconnectionresource"
```


### Client Initialization

```go
client := projectconnectionresource.NewProjectConnectionResourceClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ProjectConnectionResourceClient.ProjectConnectionsCreate`

```go
ctx := context.TODO()
id := projectconnectionresource.NewProjectConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountName", "projectName", "connectionName")

payload := projectconnectionresource.ConnectionPropertiesV2BasicResource{
	// ...
}


read, err := client.ProjectConnectionsCreate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProjectConnectionResourceClient.ProjectConnectionsDelete`

```go
ctx := context.TODO()
id := projectconnectionresource.NewProjectConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountName", "projectName", "connectionName")

read, err := client.ProjectConnectionsDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProjectConnectionResourceClient.ProjectConnectionsGet`

```go
ctx := context.TODO()
id := projectconnectionresource.NewProjectConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountName", "projectName", "connectionName")

read, err := client.ProjectConnectionsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProjectConnectionResourceClient.ProjectConnectionsList`

```go
ctx := context.TODO()
id := projectconnectionresource.NewProjectID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountName", "projectName")

// alternatively `client.ProjectConnectionsList(ctx, id, projectconnectionresource.DefaultProjectConnectionsListOperationOptions())` can be used to do batched pagination
items, err := client.ProjectConnectionsListComplete(ctx, id, projectconnectionresource.DefaultProjectConnectionsListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ProjectConnectionResourceClient.ProjectConnectionsUpdate`

```go
ctx := context.TODO()
id := projectconnectionresource.NewProjectConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountName", "projectName", "connectionName")

payload := projectconnectionresource.ConnectionUpdateContent{
	// ...
}


read, err := client.ProjectConnectionsUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

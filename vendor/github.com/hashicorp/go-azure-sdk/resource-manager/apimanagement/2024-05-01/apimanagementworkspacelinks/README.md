
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apimanagementworkspacelinks` Documentation

The `apimanagementworkspacelinks` SDK allows for interaction with Azure Resource Manager `apimanagement` (API Version `2024-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apimanagementworkspacelinks"
```


### Client Initialization

```go
client := apimanagementworkspacelinks.NewApiManagementWorkspaceLinksClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ApiManagementWorkspaceLinksClient.ApiManagementWorkspaceLinkGet`

```go
ctx := context.TODO()
id := apimanagementworkspacelinks.NewWorkspaceLinkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "workspaceId")

read, err := client.ApiManagementWorkspaceLinkGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApiManagementWorkspaceLinksClient.ListByService`

```go
ctx := context.TODO()
id := apimanagementworkspacelinks.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName")

// alternatively `client.ListByService(ctx, id)` can be used to do batched pagination
items, err := client.ListByServiceComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

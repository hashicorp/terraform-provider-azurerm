
## `github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/projectskus` Documentation

The `projectskus` SDK allows for interaction with Azure Resource Manager `devcenter` (API Version `2025-02-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/projectskus"
```


### Client Initialization

```go
client := projectskus.NewProjectSKUsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ProjectSKUsClient.SkusListByProject`

```go
ctx := context.TODO()
id := projectskus.NewProjectID("12345678-1234-9876-4563-123456789012", "example-resource-group", "projectName")

// alternatively `client.SkusListByProject(ctx, id)` can be used to do batched pagination
items, err := client.SkusListByProjectComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

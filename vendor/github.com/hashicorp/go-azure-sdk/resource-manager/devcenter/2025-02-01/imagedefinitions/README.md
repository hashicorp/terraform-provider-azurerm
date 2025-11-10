
## `github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/imagedefinitions` Documentation

The `imagedefinitions` SDK allows for interaction with Azure Resource Manager `devcenter` (API Version `2025-02-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/imagedefinitions"
```


### Client Initialization

```go
client := imagedefinitions.NewImageDefinitionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ImageDefinitionsClient.ProjectCatalogImageDefinitionBuildCancel`

```go
ctx := context.TODO()
id := imagedefinitions.NewBuildID("12345678-1234-9876-4563-123456789012", "example-resource-group", "projectName", "catalogName", "imageDefinitionName", "buildName")

if err := client.ProjectCatalogImageDefinitionBuildCancelThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ImageDefinitionsClient.ProjectCatalogImageDefinitionBuildGet`

```go
ctx := context.TODO()
id := imagedefinitions.NewBuildID("12345678-1234-9876-4563-123456789012", "example-resource-group", "projectName", "catalogName", "imageDefinitionName", "buildName")

read, err := client.ProjectCatalogImageDefinitionBuildGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ImageDefinitionsClient.ProjectCatalogImageDefinitionBuildGetBuildDetails`

```go
ctx := context.TODO()
id := imagedefinitions.NewBuildID("12345678-1234-9876-4563-123456789012", "example-resource-group", "projectName", "catalogName", "imageDefinitionName", "buildName")

read, err := client.ProjectCatalogImageDefinitionBuildGetBuildDetails(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ImageDefinitionsClient.ProjectCatalogImageDefinitionBuildsListByImageDefinition`

```go
ctx := context.TODO()
id := imagedefinitions.NewImageDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "projectName", "catalogName", "imageDefinitionName")

// alternatively `client.ProjectCatalogImageDefinitionBuildsListByImageDefinition(ctx, id)` can be used to do batched pagination
items, err := client.ProjectCatalogImageDefinitionBuildsListByImageDefinitionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ImageDefinitionsClient.ProjectCatalogImageDefinitionsBuildImage`

```go
ctx := context.TODO()
id := imagedefinitions.NewImageDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "projectName", "catalogName", "imageDefinitionName")

if err := client.ProjectCatalogImageDefinitionsBuildImageThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ImageDefinitionsClient.ProjectCatalogImageDefinitionsGetByProjectCatalog`

```go
ctx := context.TODO()
id := imagedefinitions.NewImageDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "projectName", "catalogName", "imageDefinitionName")

read, err := client.ProjectCatalogImageDefinitionsGetByProjectCatalog(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ImageDefinitionsClient.ProjectCatalogImageDefinitionsGetErrorDetails`

```go
ctx := context.TODO()
id := imagedefinitions.NewImageDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "projectName", "catalogName", "imageDefinitionName")

read, err := client.ProjectCatalogImageDefinitionsGetErrorDetails(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ImageDefinitionsClient.ProjectCatalogImageDefinitionsListByProjectCatalog`

```go
ctx := context.TODO()
id := imagedefinitions.NewCatalogID("12345678-1234-9876-4563-123456789012", "example-resource-group", "projectName", "catalogName")

// alternatively `client.ProjectCatalogImageDefinitionsListByProjectCatalog(ctx, id)` can be used to do batched pagination
items, err := client.ProjectCatalogImageDefinitionsListByProjectCatalogComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

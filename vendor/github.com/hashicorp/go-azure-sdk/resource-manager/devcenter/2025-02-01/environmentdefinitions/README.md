
## `github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/environmentdefinitions` Documentation

The `environmentdefinitions` SDK allows for interaction with Azure Resource Manager `devcenter` (API Version `2025-02-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/environmentdefinitions"
```


### Client Initialization

```go
client := environmentdefinitions.NewEnvironmentDefinitionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `EnvironmentDefinitionsClient.EnvironmentDefinitionsGet`

```go
ctx := context.TODO()
id := environmentdefinitions.NewCatalogEnvironmentDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "devCenterName", "catalogName", "environmentDefinitionName")

read, err := client.EnvironmentDefinitionsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EnvironmentDefinitionsClient.EnvironmentDefinitionsGetByProjectCatalog`

```go
ctx := context.TODO()
id := environmentdefinitions.NewEnvironmentDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "projectName", "catalogName", "environmentDefinitionName")

read, err := client.EnvironmentDefinitionsGetByProjectCatalog(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EnvironmentDefinitionsClient.EnvironmentDefinitionsGetErrorDetails`

```go
ctx := context.TODO()
id := environmentdefinitions.NewCatalogEnvironmentDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "devCenterName", "catalogName", "environmentDefinitionName")

read, err := client.EnvironmentDefinitionsGetErrorDetails(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EnvironmentDefinitionsClient.EnvironmentDefinitionsListByCatalog`

```go
ctx := context.TODO()
id := environmentdefinitions.NewDevCenterCatalogID("12345678-1234-9876-4563-123456789012", "example-resource-group", "devCenterName", "catalogName")

// alternatively `client.EnvironmentDefinitionsListByCatalog(ctx, id, environmentdefinitions.DefaultEnvironmentDefinitionsListByCatalogOperationOptions())` can be used to do batched pagination
items, err := client.EnvironmentDefinitionsListByCatalogComplete(ctx, id, environmentdefinitions.DefaultEnvironmentDefinitionsListByCatalogOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `EnvironmentDefinitionsClient.EnvironmentDefinitionsListByProjectCatalog`

```go
ctx := context.TODO()
id := environmentdefinitions.NewCatalogID("12345678-1234-9876-4563-123456789012", "example-resource-group", "projectName", "catalogName")

// alternatively `client.EnvironmentDefinitionsListByProjectCatalog(ctx, id)` can be used to do batched pagination
items, err := client.EnvironmentDefinitionsListByProjectCatalogComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `EnvironmentDefinitionsClient.ProjectCatalogEnvironmentDefinitionsGetErrorDetails`

```go
ctx := context.TODO()
id := environmentdefinitions.NewEnvironmentDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "projectName", "catalogName", "environmentDefinitionName")

read, err := client.ProjectCatalogEnvironmentDefinitionsGetErrorDetails(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```

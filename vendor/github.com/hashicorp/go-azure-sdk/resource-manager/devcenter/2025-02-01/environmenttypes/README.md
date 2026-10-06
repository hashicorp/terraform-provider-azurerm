
## `github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/environmenttypes` Documentation

The `environmenttypes` SDK allows for interaction with Azure Resource Manager `devcenter` (API Version `2025-02-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/environmenttypes"
```


### Client Initialization

```go
client := environmenttypes.NewEnvironmentTypesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `EnvironmentTypesClient.EnvironmentTypesCreateOrUpdate`

```go
ctx := context.TODO()
id := environmenttypes.NewDevCenterEnvironmentTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "devCenterName", "environmentTypeName")

payload := environmenttypes.EnvironmentType{
	// ...
}


read, err := client.EnvironmentTypesCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EnvironmentTypesClient.EnvironmentTypesDelete`

```go
ctx := context.TODO()
id := environmenttypes.NewDevCenterEnvironmentTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "devCenterName", "environmentTypeName")

read, err := client.EnvironmentTypesDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EnvironmentTypesClient.EnvironmentTypesGet`

```go
ctx := context.TODO()
id := environmenttypes.NewDevCenterEnvironmentTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "devCenterName", "environmentTypeName")

read, err := client.EnvironmentTypesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EnvironmentTypesClient.EnvironmentTypesListByDevCenter`

```go
ctx := context.TODO()
id := environmenttypes.NewDevCenterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "devCenterName")

// alternatively `client.EnvironmentTypesListByDevCenter(ctx, id, environmenttypes.DefaultEnvironmentTypesListByDevCenterOperationOptions())` can be used to do batched pagination
items, err := client.EnvironmentTypesListByDevCenterComplete(ctx, id, environmenttypes.DefaultEnvironmentTypesListByDevCenterOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `EnvironmentTypesClient.EnvironmentTypesUpdate`

```go
ctx := context.TODO()
id := environmenttypes.NewDevCenterEnvironmentTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "devCenterName", "environmentTypeName")

payload := environmenttypes.EnvironmentTypeUpdate{
	// ...
}


read, err := client.EnvironmentTypesUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EnvironmentTypesClient.ProjectAllowedEnvironmentTypesGet`

```go
ctx := context.TODO()
id := environmenttypes.NewAllowedEnvironmentTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "projectName", "allowedEnvironmentTypeName")

read, err := client.ProjectAllowedEnvironmentTypesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EnvironmentTypesClient.ProjectAllowedEnvironmentTypesList`

```go
ctx := context.TODO()
id := environmenttypes.NewProjectID("12345678-1234-9876-4563-123456789012", "example-resource-group", "projectName")

// alternatively `client.ProjectAllowedEnvironmentTypesList(ctx, id, environmenttypes.DefaultProjectAllowedEnvironmentTypesListOperationOptions())` can be used to do batched pagination
items, err := client.ProjectAllowedEnvironmentTypesListComplete(ctx, id, environmenttypes.DefaultProjectAllowedEnvironmentTypesListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `EnvironmentTypesClient.ProjectEnvironmentTypesCreateOrUpdate`

```go
ctx := context.TODO()
id := environmenttypes.NewEnvironmentTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "projectName", "environmentTypeName")

payload := environmenttypes.ProjectEnvironmentType{
	// ...
}


read, err := client.ProjectEnvironmentTypesCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EnvironmentTypesClient.ProjectEnvironmentTypesDelete`

```go
ctx := context.TODO()
id := environmenttypes.NewEnvironmentTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "projectName", "environmentTypeName")

read, err := client.ProjectEnvironmentTypesDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EnvironmentTypesClient.ProjectEnvironmentTypesGet`

```go
ctx := context.TODO()
id := environmenttypes.NewEnvironmentTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "projectName", "environmentTypeName")

read, err := client.ProjectEnvironmentTypesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EnvironmentTypesClient.ProjectEnvironmentTypesList`

```go
ctx := context.TODO()
id := environmenttypes.NewProjectID("12345678-1234-9876-4563-123456789012", "example-resource-group", "projectName")

// alternatively `client.ProjectEnvironmentTypesList(ctx, id, environmenttypes.DefaultProjectEnvironmentTypesListOperationOptions())` can be used to do batched pagination
items, err := client.ProjectEnvironmentTypesListComplete(ctx, id, environmenttypes.DefaultProjectEnvironmentTypesListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `EnvironmentTypesClient.ProjectEnvironmentTypesUpdate`

```go
ctx := context.TODO()
id := environmenttypes.NewEnvironmentTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "projectName", "environmentTypeName")

payload := environmenttypes.ProjectEnvironmentTypeUpdate{
	// ...
}


read, err := client.ProjectEnvironmentTypesUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


## `github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2023-04-01/environmentversion` Documentation

The `environmentversion` SDK allows for interaction with the Azure Resource Manager Service `machinelearningservices` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2023-04-01/environmentversion"
```


### Client Initialization

```go
client := environmentversion.NewEnvironmentVersionClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `EnvironmentVersionClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := environmentversion.NewEnvironmentVersionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "environmentValue", "versionValue")

payload := environmentversion.EnvironmentVersionResource{
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


### Example Usage: `EnvironmentVersionClient.Delete`

```go
ctx := context.TODO()
id := environmentversion.NewEnvironmentVersionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "environmentValue", "versionValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EnvironmentVersionClient.Get`

```go
ctx := context.TODO()
id := environmentversion.NewEnvironmentVersionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "environmentValue", "versionValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EnvironmentVersionClient.List`

```go
ctx := context.TODO()
id := environmentversion.NewEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "environmentValue")

// alternatively `client.List(ctx, id, environmentversion.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, environmentversion.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `EnvironmentVersionClient.RegistryEnvironmentVersionsCreateOrUpdate`

```go
ctx := context.TODO()
id := environmentversion.NewRegistryEnvironmentVersionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue", "environmentValue", "versionValue")

payload := environmentversion.EnvironmentVersionResource{
	// ...
}


if err := client.RegistryEnvironmentVersionsCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `EnvironmentVersionClient.RegistryEnvironmentVersionsDelete`

```go
ctx := context.TODO()
id := environmentversion.NewRegistryEnvironmentVersionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue", "environmentValue", "versionValue")

if err := client.RegistryEnvironmentVersionsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `EnvironmentVersionClient.RegistryEnvironmentVersionsGet`

```go
ctx := context.TODO()
id := environmentversion.NewRegistryEnvironmentVersionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue", "environmentValue", "versionValue")

read, err := client.RegistryEnvironmentVersionsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EnvironmentVersionClient.RegistryEnvironmentVersionsList`

```go
ctx := context.TODO()
id := environmentversion.NewRegistryEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue", "environmentValue")

// alternatively `client.RegistryEnvironmentVersionsList(ctx, id, environmentversion.DefaultRegistryEnvironmentVersionsListOperationOptions())` can be used to do batched pagination
items, err := client.RegistryEnvironmentVersionsListComplete(ctx, id, environmentversion.DefaultRegistryEnvironmentVersionsListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

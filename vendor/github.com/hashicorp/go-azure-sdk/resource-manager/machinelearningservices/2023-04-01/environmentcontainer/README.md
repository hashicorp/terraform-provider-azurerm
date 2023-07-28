
## `github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2023-04-01/environmentcontainer` Documentation

The `environmentcontainer` SDK allows for interaction with the Azure Resource Manager Service `machinelearningservices` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2023-04-01/environmentcontainer"
```


### Client Initialization

```go
client := environmentcontainer.NewEnvironmentContainerClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `EnvironmentContainerClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := environmentcontainer.NewEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "environmentValue")

payload := environmentcontainer.EnvironmentContainerResource{
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


### Example Usage: `EnvironmentContainerClient.Delete`

```go
ctx := context.TODO()
id := environmentcontainer.NewEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "environmentValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EnvironmentContainerClient.Get`

```go
ctx := context.TODO()
id := environmentcontainer.NewEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "environmentValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EnvironmentContainerClient.List`

```go
ctx := context.TODO()
id := environmentcontainer.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue")

// alternatively `client.List(ctx, id, environmentcontainer.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, environmentcontainer.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `EnvironmentContainerClient.RegistryEnvironmentContainersCreateOrUpdate`

```go
ctx := context.TODO()
id := environmentcontainer.NewRegistryEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue", "environmentValue")

payload := environmentcontainer.EnvironmentContainerResource{
	// ...
}


if err := client.RegistryEnvironmentContainersCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `EnvironmentContainerClient.RegistryEnvironmentContainersDelete`

```go
ctx := context.TODO()
id := environmentcontainer.NewRegistryEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue", "environmentValue")

if err := client.RegistryEnvironmentContainersDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `EnvironmentContainerClient.RegistryEnvironmentContainersGet`

```go
ctx := context.TODO()
id := environmentcontainer.NewRegistryEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue", "environmentValue")

read, err := client.RegistryEnvironmentContainersGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EnvironmentContainerClient.RegistryEnvironmentContainersList`

```go
ctx := context.TODO()
id := environmentcontainer.NewRegistryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue")

// alternatively `client.RegistryEnvironmentContainersList(ctx, id, environmentcontainer.DefaultRegistryEnvironmentContainersListOperationOptions())` can be used to do batched pagination
items, err := client.RegistryEnvironmentContainersListComplete(ctx, id, environmentcontainer.DefaultRegistryEnvironmentContainersListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

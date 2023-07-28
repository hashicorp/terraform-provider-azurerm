
## `github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2023-04-01/codecontainer` Documentation

The `codecontainer` SDK allows for interaction with the Azure Resource Manager Service `machinelearningservices` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2023-04-01/codecontainer"
```


### Client Initialization

```go
client := codecontainer.NewCodeContainerClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `CodeContainerClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := codecontainer.NewCodeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "codeValue")

payload := codecontainer.CodeContainerResource{
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


### Example Usage: `CodeContainerClient.Delete`

```go
ctx := context.TODO()
id := codecontainer.NewCodeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "codeValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CodeContainerClient.Get`

```go
ctx := context.TODO()
id := codecontainer.NewCodeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "codeValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CodeContainerClient.List`

```go
ctx := context.TODO()
id := codecontainer.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue")

// alternatively `client.List(ctx, id, codecontainer.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, codecontainer.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `CodeContainerClient.RegistryCodeContainersCreateOrUpdate`

```go
ctx := context.TODO()
id := codecontainer.NewRegistryCodeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue", "codeValue")

payload := codecontainer.CodeContainerResource{
	// ...
}


if err := client.RegistryCodeContainersCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CodeContainerClient.RegistryCodeContainersDelete`

```go
ctx := context.TODO()
id := codecontainer.NewRegistryCodeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue", "codeValue")

if err := client.RegistryCodeContainersDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CodeContainerClient.RegistryCodeContainersGet`

```go
ctx := context.TODO()
id := codecontainer.NewRegistryCodeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue", "codeValue")

read, err := client.RegistryCodeContainersGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CodeContainerClient.RegistryCodeContainersList`

```go
ctx := context.TODO()
id := codecontainer.NewRegistryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue")

// alternatively `client.RegistryCodeContainersList(ctx, id, codecontainer.DefaultRegistryCodeContainersListOperationOptions())` can be used to do batched pagination
items, err := client.RegistryCodeContainersListComplete(ctx, id, codecontainer.DefaultRegistryCodeContainersListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

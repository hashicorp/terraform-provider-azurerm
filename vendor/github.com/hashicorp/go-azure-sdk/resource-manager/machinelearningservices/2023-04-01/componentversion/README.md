
## `github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2023-04-01/componentversion` Documentation

The `componentversion` SDK allows for interaction with the Azure Resource Manager Service `machinelearningservices` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2023-04-01/componentversion"
```


### Client Initialization

```go
client := componentversion.NewComponentVersionClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ComponentVersionClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := componentversion.NewComponentVersionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "componentValue", "versionValue")

payload := componentversion.ComponentVersionResource{
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


### Example Usage: `ComponentVersionClient.Delete`

```go
ctx := context.TODO()
id := componentversion.NewComponentVersionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "componentValue", "versionValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ComponentVersionClient.Get`

```go
ctx := context.TODO()
id := componentversion.NewComponentVersionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "componentValue", "versionValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ComponentVersionClient.List`

```go
ctx := context.TODO()
id := componentversion.NewComponentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "componentValue")

// alternatively `client.List(ctx, id, componentversion.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, componentversion.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ComponentVersionClient.RegistryComponentVersionsCreateOrUpdate`

```go
ctx := context.TODO()
id := componentversion.NewRegistryComponentVersionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue", "componentValue", "versionValue")

payload := componentversion.ComponentVersionResource{
	// ...
}


if err := client.RegistryComponentVersionsCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ComponentVersionClient.RegistryComponentVersionsDelete`

```go
ctx := context.TODO()
id := componentversion.NewRegistryComponentVersionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue", "componentValue", "versionValue")

if err := client.RegistryComponentVersionsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ComponentVersionClient.RegistryComponentVersionsGet`

```go
ctx := context.TODO()
id := componentversion.NewRegistryComponentVersionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue", "componentValue", "versionValue")

read, err := client.RegistryComponentVersionsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ComponentVersionClient.RegistryComponentVersionsList`

```go
ctx := context.TODO()
id := componentversion.NewRegistryComponentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue", "componentValue")

// alternatively `client.RegistryComponentVersionsList(ctx, id, componentversion.DefaultRegistryComponentVersionsListOperationOptions())` can be used to do batched pagination
items, err := client.RegistryComponentVersionsListComplete(ctx, id, componentversion.DefaultRegistryComponentVersionsListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

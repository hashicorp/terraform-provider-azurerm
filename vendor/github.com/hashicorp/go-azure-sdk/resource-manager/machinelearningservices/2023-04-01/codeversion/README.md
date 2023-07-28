
## `github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2023-04-01/codeversion` Documentation

The `codeversion` SDK allows for interaction with the Azure Resource Manager Service `machinelearningservices` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2023-04-01/codeversion"
```


### Client Initialization

```go
client := codeversion.NewCodeVersionClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `CodeVersionClient.CreateOrGetStartPendingUpload`

```go
ctx := context.TODO()
id := codeversion.NewCodeVersionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "codeValue", "versionValue")

payload := codeversion.PendingUploadRequestDto{
	// ...
}


read, err := client.CreateOrGetStartPendingUpload(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CodeVersionClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := codeversion.NewCodeVersionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "codeValue", "versionValue")

payload := codeversion.CodeVersionResource{
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


### Example Usage: `CodeVersionClient.Delete`

```go
ctx := context.TODO()
id := codeversion.NewCodeVersionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "codeValue", "versionValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CodeVersionClient.Get`

```go
ctx := context.TODO()
id := codeversion.NewCodeVersionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "codeValue", "versionValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CodeVersionClient.List`

```go
ctx := context.TODO()
id := codeversion.NewCodeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "codeValue")

// alternatively `client.List(ctx, id, codeversion.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, codeversion.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `CodeVersionClient.RegistryCodeVersionsCreateOrGetStartPendingUpload`

```go
ctx := context.TODO()
id := codeversion.NewRegistryCodeVersionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue", "codeValue", "versionValue")

payload := codeversion.PendingUploadRequestDto{
	// ...
}


read, err := client.RegistryCodeVersionsCreateOrGetStartPendingUpload(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CodeVersionClient.RegistryCodeVersionsCreateOrUpdate`

```go
ctx := context.TODO()
id := codeversion.NewRegistryCodeVersionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue", "codeValue", "versionValue")

payload := codeversion.CodeVersionResource{
	// ...
}


if err := client.RegistryCodeVersionsCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `CodeVersionClient.RegistryCodeVersionsDelete`

```go
ctx := context.TODO()
id := codeversion.NewRegistryCodeVersionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue", "codeValue", "versionValue")

if err := client.RegistryCodeVersionsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `CodeVersionClient.RegistryCodeVersionsGet`

```go
ctx := context.TODO()
id := codeversion.NewRegistryCodeVersionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue", "codeValue", "versionValue")

read, err := client.RegistryCodeVersionsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `CodeVersionClient.RegistryCodeVersionsList`

```go
ctx := context.TODO()
id := codeversion.NewRegistryCodeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "registryValue", "codeValue")

// alternatively `client.RegistryCodeVersionsList(ctx, id, codeversion.DefaultRegistryCodeVersionsListOperationOptions())` can be used to do batched pagination
items, err := client.RegistryCodeVersionsListComplete(ctx, id, codeversion.DefaultRegistryCodeVersionsListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```

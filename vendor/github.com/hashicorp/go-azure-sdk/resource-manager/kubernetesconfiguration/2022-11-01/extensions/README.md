
## `github.com/hashicorp/go-azure-sdk/resource-manager/kubernetesconfiguration/2022-11-01/extensions` Documentation

The `extensions` SDK allows for interaction with the Azure Resource Manager Service `kubernetesconfiguration` (API Version `2022-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/kubernetesconfiguration/2022-11-01/extensions"
```


### Client Initialization

```go
client := extensions.NewExtensionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ExtensionsClient.Create`

```go
ctx := context.TODO()
id := extensions.NewScopedExtensionID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "extensionValue")

payload := extensions.Extension{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ExtensionsClient.Delete`

```go
ctx := context.TODO()
id := extensions.NewScopedExtensionID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "extensionValue")

if err := client.DeleteThenPoll(ctx, id, extensions.DefaultDeleteOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `ExtensionsClient.Get`

```go
ctx := context.TODO()
id := extensions.NewScopedExtensionID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "extensionValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ExtensionsClient.List`

```go
ctx := context.TODO()
id := extensions.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ExtensionsClient.Update`

```go
ctx := context.TODO()
id := extensions.NewScopedExtensionID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "extensionValue")

payload := extensions.PatchExtension{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

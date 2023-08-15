
## `github.com/hashicorp/go-azure-sdk/resource-manager/kubernetesconfiguration/2022-11-01/fluxconfiguration` Documentation

The `fluxconfiguration` SDK allows for interaction with the Azure Resource Manager Service `kubernetesconfiguration` (API Version `2022-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/kubernetesconfiguration/2022-11-01/fluxconfiguration"
```


### Client Initialization

```go
client := fluxconfiguration.NewFluxConfigurationClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `FluxConfigurationClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := fluxconfiguration.NewScopedFluxConfigurationID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "fluxConfigurationValue")

payload := fluxconfiguration.FluxConfiguration{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `FluxConfigurationClient.Delete`

```go
ctx := context.TODO()
id := fluxconfiguration.NewScopedFluxConfigurationID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "fluxConfigurationValue")

if err := client.DeleteThenPoll(ctx, id, fluxconfiguration.DefaultDeleteOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `FluxConfigurationClient.Get`

```go
ctx := context.TODO()
id := fluxconfiguration.NewScopedFluxConfigurationID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "fluxConfigurationValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FluxConfigurationClient.List`

```go
ctx := context.TODO()
id := fluxconfiguration.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `FluxConfigurationClient.Update`

```go
ctx := context.TODO()
id := fluxconfiguration.NewScopedFluxConfigurationID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "fluxConfigurationValue")

payload := fluxconfiguration.FluxConfigurationPatch{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

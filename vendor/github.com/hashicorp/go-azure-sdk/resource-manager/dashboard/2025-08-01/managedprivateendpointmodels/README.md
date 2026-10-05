
## `github.com/hashicorp/go-azure-sdk/resource-manager/dashboard/2025-08-01/managedprivateendpointmodels` Documentation

The `managedprivateendpointmodels` SDK allows for interaction with Azure Resource Manager `dashboard` (API Version `2025-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/dashboard/2025-08-01/managedprivateendpointmodels"
```


### Client Initialization

```go
client := managedprivateendpointmodels.NewManagedPrivateEndpointModelsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ManagedPrivateEndpointModelsClient.ManagedPrivateEndpointsCreate`

```go
ctx := context.TODO()
id := managedprivateendpointmodels.NewManagedPrivateEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "grafanaName", "managedPrivateEndpointName")

payload := managedprivateendpointmodels.ManagedPrivateEndpointModel{
	// ...
}


if err := client.ManagedPrivateEndpointsCreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ManagedPrivateEndpointModelsClient.ManagedPrivateEndpointsDelete`

```go
ctx := context.TODO()
id := managedprivateendpointmodels.NewManagedPrivateEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "grafanaName", "managedPrivateEndpointName")

if err := client.ManagedPrivateEndpointsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ManagedPrivateEndpointModelsClient.ManagedPrivateEndpointsGet`

```go
ctx := context.TODO()
id := managedprivateendpointmodels.NewManagedPrivateEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "grafanaName", "managedPrivateEndpointName")

read, err := client.ManagedPrivateEndpointsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedPrivateEndpointModelsClient.ManagedPrivateEndpointsList`

```go
ctx := context.TODO()
id := managedprivateendpointmodels.NewGrafanaID("12345678-1234-9876-4563-123456789012", "example-resource-group", "grafanaName")

// alternatively `client.ManagedPrivateEndpointsList(ctx, id)` can be used to do batched pagination
items, err := client.ManagedPrivateEndpointsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ManagedPrivateEndpointModelsClient.ManagedPrivateEndpointsUpdate`

```go
ctx := context.TODO()
id := managedprivateendpointmodels.NewManagedPrivateEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "grafanaName", "managedPrivateEndpointName")

payload := managedprivateendpointmodels.ManagedPrivateEndpointUpdateParameters{
	// ...
}


if err := client.ManagedPrivateEndpointsUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

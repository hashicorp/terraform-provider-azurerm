
## `github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/providerinstances` Documentation

The `providerinstances` SDK allows for interaction with Azure Resource Manager `workloads` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/providerinstances"
```


### Client Initialization

```go
client := providerinstances.NewProviderInstancesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ProviderInstancesClient.Create`

```go
ctx := context.TODO()
id := providerinstances.NewProviderInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName", "providerInstanceName")

payload := providerinstances.ProviderInstance{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ProviderInstancesClient.Delete`

```go
ctx := context.TODO()
id := providerinstances.NewProviderInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName", "providerInstanceName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ProviderInstancesClient.Get`

```go
ctx := context.TODO()
id := providerinstances.NewProviderInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName", "providerInstanceName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ProviderInstancesClient.List`

```go
ctx := context.TODO()
id := providerinstances.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


## `github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/integrationruntimes` Documentation

The `integrationruntimes` SDK allows for interaction with Azure Resource Manager `datafactory` (API Version `2018-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/integrationruntimes"
```


### Client Initialization

```go
client := integrationruntimes.NewIntegrationRuntimesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `IntegrationRuntimesClient.CreateLinkedIntegrationRuntime`

```go
ctx := context.TODO()
id := integrationruntimes.NewIntegrationRuntimeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "factoryName", "integrationRuntimeName")

payload := integrationruntimes.CreateLinkedIntegrationRuntimeRequest{
	// ...
}


read, err := client.CreateLinkedIntegrationRuntime(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IntegrationRuntimesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := integrationruntimes.NewIntegrationRuntimeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "factoryName", "integrationRuntimeName")

payload := integrationruntimes.IntegrationRuntimeResource{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, integrationruntimes.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IntegrationRuntimesClient.Delete`

```go
ctx := context.TODO()
id := integrationruntimes.NewIntegrationRuntimeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "factoryName", "integrationRuntimeName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IntegrationRuntimesClient.Get`

```go
ctx := context.TODO()
id := integrationruntimes.NewIntegrationRuntimeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "factoryName", "integrationRuntimeName")

read, err := client.Get(ctx, id, integrationruntimes.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IntegrationRuntimesClient.GetConnectionInfo`

```go
ctx := context.TODO()
id := integrationruntimes.NewIntegrationRuntimeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "factoryName", "integrationRuntimeName")

read, err := client.GetConnectionInfo(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IntegrationRuntimesClient.GetMonitoringData`

```go
ctx := context.TODO()
id := integrationruntimes.NewIntegrationRuntimeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "factoryName", "integrationRuntimeName")

read, err := client.GetMonitoringData(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IntegrationRuntimesClient.GetStatus`

```go
ctx := context.TODO()
id := integrationruntimes.NewIntegrationRuntimeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "factoryName", "integrationRuntimeName")

read, err := client.GetStatus(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IntegrationRuntimesClient.ListAuthKeys`

```go
ctx := context.TODO()
id := integrationruntimes.NewIntegrationRuntimeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "factoryName", "integrationRuntimeName")

read, err := client.ListAuthKeys(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IntegrationRuntimesClient.ListByFactory`

```go
ctx := context.TODO()
id := integrationruntimes.NewFactoryID("12345678-1234-9876-4563-123456789012", "example-resource-group", "factoryName")

// alternatively `client.ListByFactory(ctx, id)` can be used to do batched pagination
items, err := client.ListByFactoryComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `IntegrationRuntimesClient.ListOutboundNetworkDependenciesEndpoints`

```go
ctx := context.TODO()
id := integrationruntimes.NewIntegrationRuntimeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "factoryName", "integrationRuntimeName")

read, err := client.ListOutboundNetworkDependenciesEndpoints(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IntegrationRuntimesClient.RegenerateAuthKey`

```go
ctx := context.TODO()
id := integrationruntimes.NewIntegrationRuntimeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "factoryName", "integrationRuntimeName")

payload := integrationruntimes.IntegrationRuntimeRegenerateKeyParameters{
	// ...
}


read, err := client.RegenerateAuthKey(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IntegrationRuntimesClient.RemoveLinks`

```go
ctx := context.TODO()
id := integrationruntimes.NewIntegrationRuntimeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "factoryName", "integrationRuntimeName")

payload := integrationruntimes.LinkedIntegrationRuntimeRequest{
	// ...
}


read, err := client.RemoveLinks(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IntegrationRuntimesClient.Start`

```go
ctx := context.TODO()
id := integrationruntimes.NewIntegrationRuntimeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "factoryName", "integrationRuntimeName")

if err := client.StartThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `IntegrationRuntimesClient.Stop`

```go
ctx := context.TODO()
id := integrationruntimes.NewIntegrationRuntimeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "factoryName", "integrationRuntimeName")

if err := client.StopThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `IntegrationRuntimesClient.SyncCredentials`

```go
ctx := context.TODO()
id := integrationruntimes.NewIntegrationRuntimeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "factoryName", "integrationRuntimeName")

read, err := client.SyncCredentials(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IntegrationRuntimesClient.Update`

```go
ctx := context.TODO()
id := integrationruntimes.NewIntegrationRuntimeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "factoryName", "integrationRuntimeName")

payload := integrationruntimes.UpdateIntegrationRuntimeRequest{
	// ...
}


read, err := client.Update(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IntegrationRuntimesClient.Upgrade`

```go
ctx := context.TODO()
id := integrationruntimes.NewIntegrationRuntimeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "factoryName", "integrationRuntimeName")

read, err := client.Upgrade(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


## `github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2023-04-01/operationalizationclusters` Documentation

The `operationalizationclusters` SDK allows for interaction with the Azure Resource Manager Service `machinelearningservices` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2023-04-01/operationalizationclusters"
```


### Client Initialization

```go
client := operationalizationclusters.NewOperationalizationClustersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `OperationalizationClustersClient.ComputeCreateOrUpdate`

```go
ctx := context.TODO()
id := operationalizationclusters.NewComputeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "computeValue")

payload := operationalizationclusters.ComputeResource{
	// ...
}


if err := client.ComputeCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OperationalizationClustersClient.ComputeDelete`

```go
ctx := context.TODO()
id := operationalizationclusters.NewComputeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "computeValue")

if err := client.ComputeDeleteThenPoll(ctx, id, operationalizationclusters.DefaultComputeDeleteOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `OperationalizationClustersClient.ComputeGet`

```go
ctx := context.TODO()
id := operationalizationclusters.NewComputeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "computeValue")

read, err := client.ComputeGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OperationalizationClustersClient.ComputeList`

```go
ctx := context.TODO()
id := operationalizationclusters.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue")

// alternatively `client.ComputeList(ctx, id, operationalizationclusters.DefaultComputeListOperationOptions())` can be used to do batched pagination
items, err := client.ComputeListComplete(ctx, id, operationalizationclusters.DefaultComputeListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OperationalizationClustersClient.ComputeListKeys`

```go
ctx := context.TODO()
id := operationalizationclusters.NewComputeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "computeValue")

read, err := client.ComputeListKeys(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OperationalizationClustersClient.ComputeRestart`

```go
ctx := context.TODO()
id := operationalizationclusters.NewComputeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "computeValue")

if err := client.ComputeRestartThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OperationalizationClustersClient.ComputeStart`

```go
ctx := context.TODO()
id := operationalizationclusters.NewComputeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "computeValue")

if err := client.ComputeStartThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OperationalizationClustersClient.ComputeStop`

```go
ctx := context.TODO()
id := operationalizationclusters.NewComputeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "computeValue")

if err := client.ComputeStopThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OperationalizationClustersClient.ComputeUpdate`

```go
ctx := context.TODO()
id := operationalizationclusters.NewComputeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "computeValue")

payload := operationalizationclusters.ClusterUpdateParameters{
	// ...
}


if err := client.ComputeUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


## `github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2023-04-01/batchdeployment` Documentation

The `batchdeployment` SDK allows for interaction with the Azure Resource Manager Service `machinelearningservices` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2023-04-01/batchdeployment"
```


### Client Initialization

```go
client := batchdeployment.NewBatchDeploymentClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `BatchDeploymentClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := batchdeployment.NewDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "batchEndpointValue", "deploymentValue")

payload := batchdeployment.BatchDeploymentTrackedResource{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `BatchDeploymentClient.Delete`

```go
ctx := context.TODO()
id := batchdeployment.NewDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "batchEndpointValue", "deploymentValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `BatchDeploymentClient.Get`

```go
ctx := context.TODO()
id := batchdeployment.NewDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "batchEndpointValue", "deploymentValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BatchDeploymentClient.List`

```go
ctx := context.TODO()
id := batchdeployment.NewBatchEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "batchEndpointValue")

// alternatively `client.List(ctx, id, batchdeployment.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, batchdeployment.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `BatchDeploymentClient.Update`

```go
ctx := context.TODO()
id := batchdeployment.NewDeploymentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceValue", "batchEndpointValue", "deploymentValue")

payload := batchdeployment.PartialBatchDeploymentPartialMinimalTrackedResourceWithProperties{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```

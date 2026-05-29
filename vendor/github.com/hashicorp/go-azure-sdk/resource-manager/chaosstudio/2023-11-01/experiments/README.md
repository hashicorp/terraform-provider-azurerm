
## `github.com/hashicorp/go-azure-sdk/resource-manager/chaosstudio/2023-11-01/experiments` Documentation

The `experiments` SDK allows for interaction with Azure Resource Manager `chaosstudio` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/chaosstudio/2023-11-01/experiments"
```


### Client Initialization

```go
client := experiments.NewExperimentsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ExperimentsClient.Cancel`

```go
ctx := context.TODO()
id := experiments.NewExperimentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "experimentName")

if err := client.CancelThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ExperimentsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := experiments.NewExperimentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "experimentName")

payload := experiments.Experiment{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ExperimentsClient.Delete`

```go
ctx := context.TODO()
id := experiments.NewExperimentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "experimentName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ExperimentsClient.ExecutionDetails`

```go
ctx := context.TODO()
id := experiments.NewExecutionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "experimentName", "executionId")

read, err := client.ExecutionDetails(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ExperimentsClient.Get`

```go
ctx := context.TODO()
id := experiments.NewExperimentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "experimentName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ExperimentsClient.GetExecution`

```go
ctx := context.TODO()
id := experiments.NewExecutionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "experimentName", "executionId")

read, err := client.GetExecution(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ExperimentsClient.List`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.List(ctx, id, experiments.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, experiments.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ExperimentsClient.ListAll`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListAll(ctx, id, experiments.DefaultListAllOperationOptions())` can be used to do batched pagination
items, err := client.ListAllComplete(ctx, id, experiments.DefaultListAllOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ExperimentsClient.ListAllExecutions`

```go
ctx := context.TODO()
id := experiments.NewExperimentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "experimentName")

// alternatively `client.ListAllExecutions(ctx, id)` can be used to do batched pagination
items, err := client.ListAllExecutionsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ExperimentsClient.Start`

```go
ctx := context.TODO()
id := experiments.NewExperimentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "experimentName")

if err := client.StartThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ExperimentsClient.Update`

```go
ctx := context.TODO()
id := experiments.NewExperimentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "experimentName")

payload := experiments.ExperimentUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
